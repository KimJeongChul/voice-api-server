package api

import (
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	"voice-api-server/common"
	"voice-api-server/errors"
	"voice-api-server/logger"
	"voice-api-server/utils"
	"voice-api-server/voiceExtension"
)

// handleSpeech Serving STT service /v1/speech
func (as *ApiServer) handleSpeech(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	mw := multipart.NewWriter(w)
	var session *common.SessionObj
	var cErr *errors.CError
	var speechMetaData MsgSpeechMetaData
	errCode := http.StatusInternalServerError
	codec := "pcm"

	var sttPushedTime int
	var convertedSize, receivedSize int

	//For MultiPart, Chunked Response
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", mw.FormDataContentType())

	// Preprocessing api call
	session, cErr, errCode = as.APIPreProcessing(&w, req)
	if cErr != nil {
		if session != nil {
			logger.LogE(session.FuncName, session.TransactionId, "API Preprocessing error: ", cErr.Error())
		}
		http.Error(w, cErr.Error(), errCode)
	}
	defer func() {
		if cErr != nil {
			logger.LogE(session.FuncName, session.TransactionId, "Error:Msg=", cErr.Error())
			http.Error(w, cErr.Error(), errCode)
			responseTime := time.Now().Sub(startTime).Milliseconds()
			logger.LogI(session.FuncName, session.TransactionId, "responsetime=", responseTime)
		}
	}()

	// Parse MIME media type
	contentType, params, parseErr := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if parseErr != nil || !strings.HasPrefix(contentType, "multipart") {
		cErr, errCode = errors.NewCError(errors.HTTP_MULTIPART_PARSE_ERR, "Content-Type or multipart is not defined"), http.StatusBadRequest
		return
	}

	//Set MultiPart-Reader
	multipartReader := multipart.NewReader(req.Body, params["boundary"])
	defer req.Body.Close()

	isRecvMultiPartMetadata := false
	isRecvMultiPartMedia := false

	for {
		//Read MultiPart
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			if !isRecvMultiPartMedia || !isRecvMultiPartMetadata {
				logger.LogI(session.FuncName, session.TransactionId, "multipart reader io EOF error :", err.Error())
				cErr, errCode = errors.NewCError(errors.HTTP_MULTIPART_READ_ERR, "multipart reader error : "+err.Error()+", Please check multipart name: metadata, media"), http.StatusInternalServerError
			}
			break
		}
		if err != nil {
			cErr, errCode = errors.NewCError(errors.HTTP_MULTIPART_READ_ERR, err.Error()), http.StatusInternalServerError
			logger.LogI(session.FuncName, session.TransactionId, "multipart reader error:", cErr.Error())
			return
		}
		defer part.Close()

		// Read Part
		fileBytes, err := ioutil.ReadAll(part)
		if err != nil {
			cErr, errCode = errors.NewCError(errors.IOUTIL_READ_ERR, err.Error()), http.StatusInternalServerError
			return
		}
		fileBytesLength := len(fileBytes)


		switch part.FormName() {
		case "metadata":
			isRecvMultiPartMetadata = true
			err := json.Unmarshal(fileBytes, &speechMetaData)
			if err != nil {
				cErr, errCode = errors.NewCError(errors.JSON_UNMARSHAL_ERR, err.Error()), http.StatusBadRequest
				return
			}
			logger.LogI(session.FuncName, session.TransactionId, "req body: ", string(fileBytes))
			logger.LogI(session.FuncName, session.TransactionId, "Read Meta:Encoding="+speechMetaData.Encoding+",Channel="+strconv.Itoa(speechMetaData.EncodingOpt.Channel)+",sampleRate="+strconv.Itoa(speechMetaData.EncodingOpt.SampleRate)+",targetLanguage="+speechMetaData.TargetLanguage)
			//Check Necessary Parameters
			cErr, errCode = checkSpeechReqParams(speechMetaData), http.StatusBadRequest
			if cErr != nil {
				return
			}
		case "media":
			isRecvMultiPartMedia = true
			logger.LogI(session.FuncName, session.TransactionId, "Read Audio:Length="+strconv.Itoa(fileBytesLength))
		
			receivedSize = fileBytesLength

			maxBufferLength := 16000 * 2 * 60 * 10
			var maxAudioBuffer []byte
			var inPcmBuffer []byte

			finalResultMsg := FileSpeechRecognizeRes{SttStatus: "completed", TransactionId: session.TransactionId}
			finalSttResult := make([]MsgSpeechRecognizeResSttResult, 0)
			finalResultMsg.SttResults = &finalSttResult


			// Check audio file size is zero
			if receivedSize <= 0 {
				cErr, errCode = errors.NewCError(errors.HTTP_BODY_READ_ERR, "Receive audio file size is zero"), http.StatusOK
				return
			}

			if !isRecvMultiPartMetadata {
				cErr, errCode = errors.NewCError(errors.HTTP_MULTIPART_READ_ERR, "multipart reader error : Please check multipart name: metadata"), http.StatusBadRequest
				return
			}

			// Save audio file 
			saveFileName := utils.GetFileSaveFormattedTime() + "." + speechMetaData.Encoding
			saveFilePath := as.serverConfig.RcvAudioSavePath + "/" + saveFileName
			ioutil.WriteFile(saveFilePath, fileBytes, 0644)
			
			logger.LogI(session.FuncName, session.TransactionId, "encoding="+speechMetaData.Encoding)

			// Convert PCM data which is stt core engin model input format
			if speechMetaData.Encoding != "raw" {
				maxAudioBuffer = make([]byte, maxBufferLength)

				if speechMetaData.Encoding != "wav" {
					if speechMetaData.Encoding == "aac" {
						// Encodig AAC(M4A)
						codec, convertedSize = voiceExtension.GetPcmFromM4aFile(saveFilePath, maxAudioBuffer)
						if convertedSize < 0 {
							cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "Audio Converting ERROR"), http.StatusOK
							return
						}
					} else {
						codec, convertedSize = voiceExtension.GetPcmFromEncoded(fileBytes, maxAudioBuffer)
						if convertedSize < 0 {
							cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "Audio Converting ERROR"), http.StatusOK
							return
						}
					}
				} else {

					// Encoding WAV
					codec = "wav"
					convertedSize = voiceExtension.GetWavToPcm(fileBytes, maxAudioBuffer)
					if convertedSize < 0 {
						cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "Audio Converting ERROR"), http.StatusOK
						return
					}	
				}
				inPcmBuffer = maxAudioBuffer[0:convertedSize]
			}  else {
				// Encoding Raw PCM
				if speechMetaData.EncodingOpt.SampleFmt != "S16LE" || speechMetaData.EncodingOpt.SampleRate != 16000 || speechMetaData.EncodingOpt.Channel != 1 {
					maxAudioBuffer = make([]byte, maxBufferLength)
					//codec, convertedSize = voiceExtension.GetResampledPcm(speechMetaData.EncodingOpt.SampleFmt, speechMetaData.EncodingOpt.SampleRate, speechMetaData.EncodingOpt.Channel, fileBytes, maxAudioBuffer)
					logger.LogI(session.FuncName, session.TransactionId, "Resampling:sourceDataLength=", fileBytesLength, ",:destinationDataLength=", convertedSize)
					if convertedSize < 0 {
						if convertedSize == -500 {
							cErr, errCode = errors.NewCError(errors.STT_RESAMPLE_ERR, "cannot initialize RESAMPLER"), http.StatusInternalServerError
						} else {
							cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "cannot resample data"), http.StatusOK
						}
						return
					}
					inPcmBuffer = maxAudioBuffer[0:convertedSize]
				} else {
					if fileBytesLength > maxBufferLength {
						inPcmBuffer = fileBytes[0:maxBufferLength]
						convertedSize = maxBufferLength
					} else {
						inPcmBuffer = fileBytes[0:fileBytesLength]
						convertedSize = fileBytesLength
					}
				}
			}
			logger.LogI(session.FuncName, session.TransactionId, "codec=", codec)
			logger.LogI(session.FuncName, session.TransactionId, "inPcmBuffer length=", len(inPcmBuffer))

			sttPushedTime = len(inPcmBuffer) / 16000 / 2

			// Save audio pcm file 
			savePCMFileName := utils.GetFileSaveFormattedTime() + ".pcm"
			savePCMFilePath := as.serverConfig.PcmSavePath + "/" + savePCMFileName
			ioutil.WriteFile(savePCMFilePath, inPcmBuffer, 0644)

			/**
			 * RQUEST YOUR CUSTOM STT Core Engine SERVER
			*/
			// Start Recognize Voice
			msgStartRes := MsgSpeechRecognizeRes{ResultType: "start"}
			wByte, err := json.Marshal(msgStartRes)
			if err != nil {
				cErr, errCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
				return
			}
			mw.WriteField("voiceResult", string(wByte))
			session.Flush()

			// Result Recognize
			sttRecogText := "RECOGNIZE FROM YOUR STT CORE ENGIN SERVER"
			resultType := "full"
			sttStartTime := float32(0.1)
			sttEndTime := float32(40.0)

			var returnMsg MsgSpeechRecognizeRes
			returnMsg.ResultType = resultType
			returnMsg.SttResult = &MsgSpeechRecognizeResSttResult{Text: sttRecogText, StartTime: sttStartTime, EndTime: sttEndTime}
			wByte, err = json.Marshal(returnMsg)
			if err != nil {
				cErr, errCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
				return
			}
			mw.WriteField("voiceResult", string(wByte))
			session.Flush()

			finalSttResult = append(finalSttResult, *returnMsg.SttResult)

			/**
			 * END COMMUNICATION STT STT Core Engine SERVE
			*/

			var speechInfoMsg MsgSpeechRecognizeResSttInfo
			speechInfoMsg.ReqFileSize = fileBytesLength
			speechInfoMsg.TransCodec = codec
			speechInfoMsg.ConvFileSize = convertedSize
			speechInfoMsg.SttInputTime = sttPushedTime

			finalResultMsg.SttInfo = &speechInfoMsg
			finalResultMsg.EventTime = utils.GetMillisTimeFormat(startTime)

			finalMsg := MsgSpeechRecognizeRes{ResultType: "end"}
			finalMsg.SttInfo = &MsgSpeechRecognizeResSttInfo{ReqFileSize: fileBytesLength, ConvFileSize: convertedSize, SttInputTime: sttPushedTime, TransCodec: codec}
			wByte, err = json.Marshal(finalMsg)
			if err != nil {
				cErr, errCode = errors.NewCError(errors.JSON_MARSHAL_ERR, err.Error()), http.StatusInternalServerError
				return
			}
			mw.WriteField("voiceResult", string(wByte))
			session.Flush()
			mw.Close()

			// Write File Speech Recognize
			wByte, _ = json.Marshal(finalResultMsg)
			as.speechLogger.WriteByte(wByte)

			return 
		}
	}

	responseTime := time.Now().Sub(startTime).Milliseconds()
	logger.LogI(session.FuncName, session.TransactionId, "responsetime=", responseTime)
}

// handleSynthesis Serving TTS service  /v1/synthesis
func (as *ApiServer) handleSynthesis(w http.ResponseWriter, req *http.Request) {

}
