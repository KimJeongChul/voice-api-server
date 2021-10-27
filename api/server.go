package api

import (
	"net/http"

	"voice-api-server/errors"
	"voice-api-server/logger"
	"voice-api-server/utils"

	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type ApiServer struct {
	serverConfig *utils.ServerConfigJson
	router       *mux.Router
	logger       *logger.Logger
}

func (as *ApiServer) Initialize(configJson *utils.ServerConfigJson, rlogger *rotatelogs.RotateLogs) int {
	as.logger = logger.NewLogger(rlogger)
	as.serverConfig = configJson

	// Router
	as.router = mux.NewRouter()

	// STT
	as.router.HandleFunc("/v1/speech", as.handleSpeech)
	as.router.HandleFunc("/v1/synthesis", as.handleSynthesis)

	return 0
}

//Listen HTTP Serve
func (as *ApiServer) Listen() int {
	var err error
	if as.serverConfig.Ssl == 0 {
		logger.LogI("Listen", "api", "HTTP Listening:Port=", as.serverConfig.ListenPort)
		err = http.ListenAndServe(":"+as.serverConfig.ListenPort, as.router)
		cErr := errors.NewCError(errors.HTTP_SERVE_ERR, err.Error())
		logger.LogE("Listen", "api", "ERR:Msg=", cErr.Error())
	} else {
		logger.LogI("Listen", "api", "HTTPs Listening:Cert=", as.serverConfig.CertPemPath, ",keyPem=", as.serverConfig.KeyPemPath)
		err = http.ListenAndServeTLS(":"+as.serverConfig.ListenPort, as.serverConfig.CertPemPath, as.serverConfig.KeyPemPath, as.router)
		cErr := errors.NewCError(errors.HTTPS_SERVE_ERR, err.Error())
		logger.LogE("Listen", "api", "ERR:Msg=", cErr.Error())
	}
	if err != nil {
		panic(err)
	}
	return 0
}
