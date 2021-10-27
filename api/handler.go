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
	"voice-api-server/voiceExtension"
)

// handleSpeech Serving STT service /v1/speech
func (as *ApiServer) handleSpeech(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	mw := multipart.NewWriter(w)
	var session *common.SessionObj
	var cErr *errors.CError
	errCode := http.StatusInternalServerError
	codec := "pcm"

	//var sttPushedTime float32
	var convertedSize, receivedSize int
	//var receivedSize int

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

		var speechMetaData MsgSpeechMetaData

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

			// Check audio file size is zero
			if receivedSize <= 0 {
				cErr, errCode = errors.NewCError(errors.HTTP_BODY_READ_ERR, "Receive audio file size is zero"), http.StatusOK
				return
			}

			if !isRecvMultiPartMetadata {
				cErr, errCode = errors.NewCError(errors.HTTP_MULTIPART_READ_ERR, "multipart reader error : Please check multipart name: metadata"), http.StatusBadRequest
				return
			}

			// Convert PCM data which is stt core engin model input format
			if speechMetaData.Encoding != "raw" {
				maxAudioBuffer = make([]byte, maxBufferLength)
				if speechMetaData.Encoding != "wav" {
					if speechMetaData.Encoding == "aac" {
					} else {
						codec, convertedSize = voiceExtension.GetPcmFromEncoded(fileBytes, maxAudioBuffer)
						if convertedSize < 0 {
							cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "Audio Converting ERROR"), http.StatusOK
							return
						}
					}
				} else {
					// Encoding wav fiel
					codec = "wav"
					convertedSize = voiceExtension.GetWavToPcm(fileBytes, maxAudioBuffer)
					if convertedSize < 0 {
						cErr, errCode = errors.NewCError(errors.AUDIO_ENCODING_BUFFER_ERR, "Audio Converting ERROR"), http.StatusOK
						return
					}	
				}

			}  else {
			}
			logger.LogI(session.FuncName, session.TransactionId, "codec=", codec)

		}
	}

	responseTime := time.Now().Sub(startTime).Milliseconds()
	logger.LogI(session.FuncName, session.TransactionId, "responsetime=", responseTime)
}

// handleSynthesis Serving TTS service  /v1/synthesis
func (as *ApiServer) handleSynthesis(w http.ResponseWriter, req *http.Request) {

}
