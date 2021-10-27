package logger

import (
	"fmt"
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type Logger struct {
	logger *rotatelogs.RotateLogs
}

func NewLogger(logger *rotatelogs.RotateLogs) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Write(msg string) {
	l.logger.Write([]byte(l.genLogTimeStamp() + " "))
	l.logger.Write([]byte(msg))
	l.logger.Write([]byte("\n"))
}

func (l *Logger) WriteByte(msg []byte) {
	l.logger.Write(msg)
	l.logger.Write([]byte("\n"))
}

func (l *Logger) genLogTimeStamp() string {
	t := time.Now()
	return string(t.Format("2006/01/02 15:04:05"))
}

func LogE(funcName string, trxId string, logMsg ...interface{}) {
	logStr := "[E] ( " + trxId + " ) < " + funcName + " > "
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(string(colorRed), logStr, string(colorReset))
}
func LogI(funcName string, trxId string, logMsg ...interface{}) {
	logStr := string(colorBlue) + "[I] ( " + trxId + " ) < " + funcName + " > " + string(colorReset)
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}
func LogD(funcName string, trxId string, logMsg ...interface{}) {
	logStr := string(colorGreen) + "[D] ( " + trxId + " ) < " + funcName + " > " + string(colorReset)
	for _, item := range logMsg {
		logStr += fmt.Sprint(item)
	}
	log.Println(logStr)
}

func Startup() {
	// https://lunicode.com/bigtext
	fmt.Println(string(colorPurple))
	log.Println("__/\\\\\\________/\\\\\\_______________________________________________________")
	log.Println(" _\\/\\\\\\_______\\/\\\\\\_______________________________________________________")
	log.Println("  _\\//\\\\\\______/\\\\\\_________________/\\\\\\___________________________________")
	log.Println("   __\\//\\\\\\____/\\\\\\_______/\\\\\\\\\\____\\///______/\\\\\\\\\\\\\\\\_____/\\\\\\\\\\\\\\\\________")
	log.Println("    ___\\//\\\\\\__/\\\\\\______/\\\\\\///\\\\\\___/\\\\\\___/\\\\\\//////____/\\\\\\/////\\\\\\________")
	log.Println("     ____\\//\\\\\\/\\\\\\______/\\\\\\__\\//\\\\\\_\\/\\\\\\__/\\\\\\__________/\\\\\\\\\\\\\\\\\\\\\\____")
	log.Println("      _____\\//\\\\\\\\\\______\\//\\\\\\__/\\\\\\__\\/\\\\\\_\\//\\\\\\________\\//\\\\///////__________")
	log.Println("       ______\\//\\\\\\________\\///\\\\\\\\\\/___\\/\\\\\\__\\///\\\\\\\\\\\\\\\\__\\//\\\\\\\\\\\\\\\\\\\\_____")
	log.Println("        _______\\///___________\\/////_____\\///_____\\////////____\\//////////__________", string(colorReset))
}
