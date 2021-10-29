package main

import (
	"flag"
	"log"
	"os"
	"time"
	"voice-api-server/api"
	"voice-api-server/common"
	"voice-api-server/logger"
	"voice-api-server/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func main() {
	// Logger start-up
	logger.Startup()

	// Parse arguments
	configFilePath := flag.String("c", "/mnt/dev/gopath/src/voice-api-server/serverConfig.json", "Set server config file")
	flag.Parse()

	// Load server configuaration file
	config, err := utils.LoadConfigJson(configFilePath)
	if err != nil {
		logger.LogE("main", "UNDEFINED", "Config File:"+*configFilePath+" Load Error.")
		os.Exit(-1)
	}

	logPath := config.LogPath + "/%Y%m%d.debug"
	rlogger, err := rotatelogs.New(
		logPath,
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(config.LogPeriod)),
	)
	if err != nil {
		logger.LogE("main", "UNDEFINED", "Config File:"+*configFilePath+" Load Error.")
		os.Exit(-1)
	}

	logSpeechResultPath := config.SpeechResultPath + "/%Y%m%d.debug"
	rSpeechlogger, err := rotatelogs.New(
		logSpeechResultPath,
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(config.LogPeriod)),
	)
	if err != nil {
		logger.LogE("main", "UNDEFINED", "Config File:"+*configFilePath+" Load Error.")
		os.Exit(-1)
	}

	log.SetOutput(rlogger)

	logger.LogI("main", common.UNDEFINED, "Listening Config:HttpPort=", config.ListenPort)
	logger.LogI("main", common.UNDEFINED, "Extra:logLevel=", config.LogLevel, ",ssl=", config.Ssl, ",CertPemPath=", config.CertPemPath, ",KeyPemPath="+config.KeyPemPath)

	apiServer := api.ApiServer{}
	apiServer.Initialize(&config, rlogger, rSpeechlogger)


	//Listen
	apiServer.Listen()
}