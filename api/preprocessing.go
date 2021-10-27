package api

import (
	"net/http"
	"runtime"
	"strings"
	"voice-api-server/utils"
	"voice-api-server/errors"
	"voice-api-server/logger"
	"voice-api-server/common"
)

func (as *ApiServer) APIPreProcessing(w *http.ResponseWriter, req *http.Request)  (sessionObj *common.SessionObj, cErr *errors.CError, errCode int) {
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