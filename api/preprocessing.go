package api

import (
	"net/http"
	"runtime"
	"strings"
	"voice-api-server/common"
	"voice-api-server/errors"
	"voice-api-server/logger"
	"voice-api-server/utils"
)

func (as *ApiServer) APIPreProcessing(w *http.ResponseWriter, req *http.Request) (sessionObj *common.SessionObj, cErr *errors.CError, errCode int) {
	sessionObj = &common.SessionObj{}
	var err error

	utils.EnableCors(w)

	if req.Method == "OPTIONS" {
		cErr = errors.NewCError(errors.HTTP_INVALID_METHOD_ERR, "Options Request. Ignore")
		logger.LogE(sessionObj.FuncName, sessionObj.TransactionId, "ERROR:Msg=", cErr.Error())
		return
	}

	sessionObj.TransactionId, err = utils.GenerateUnixTrxID()
	if err != nil {
		cErr = errors.NewCError(errors.HTTP_PREPROCESSING_ERR, "Cannot Create TransactionId")
		return
	}

	//Create Flusher
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		cErr = errors.NewCError(errors.HTTP_FLUSHER_ERR, "Cannot Create Flush Object")
		return
	}
	sessionObj.Flush = func() {
		flusher.Flush()
	}

	//Set Function Name
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		sessionObj.FuncName = "UNKNOWN"
	} else {
		funcNameArr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		sessionObj.FuncName = funcNameArr[len(funcNameArr)-1]
	}

	return
}

// checkSpeechReqParams Check valid parameters
func checkSpeechReqParams(target MsgSpeechMetaData) (cErr *errors.CError) {
	//encoding parameter check
	if !(target.Encoding == "raw" || target.Encoding == "mp3" || target.Encoding == "vor" || target.Encoding == "wav" || target.Encoding == "aac" || target.Encoding == "fla") {
		cErr = errors.NewCError(errors.STT_FORMAT_ERR, "not supported audio format")
		return
	}
	// Check encoding is raw
	if target.Encoding == "raw" {
		//channel parameter check
		if !(target.EncodingOpt.Channel == 1 || target.EncodingOpt.Channel == 2) {
			cErr = errors.NewCError(errors.STT_CHANNEL_ERR, "not supported audio channel")
			return
		}
		//sampleRate Check
		if !(target.EncodingOpt.SampleRate == 16000 || target.EncodingOpt.SampleRate == 44100 || target.EncodingOpt.SampleRate == 48000) {
			cErr = errors.NewCError(errors.STT_SAMPLERATE_ERR, "not supported sample rate")
			return
		}
		if !(target.EncodingOpt.SampleFmt == "S16LE" || target.EncodingOpt.SampleFmt == "F32LE") {
			cErr = errors.NewCError(errors.STT_FORMAT_ERR, "not supported audio format")
			return
		}
	}
	return nil
}
