package utils

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)


var uuidCnt int
var localMacAddr string

func init() {
	localMacAddrs, err := getMacAddress()
	if err != nil {
		panic(err)
	}

	localMacAddr = localMacAddrs[0]
	uuidCnt = 0
}

type ServerConfigJson struct {
	ServiceUrl        string   `json:"serviceUrl"`
	ListenPort        string   `json:"listenPort"`
	Ssl               int      `json:"ssl"`
	CertPemPath       string   `json:"certPemPath"`
	KeyPemPath        string   `json:"keyPemPath"`
	RcvAudioSavePath  string   `json:"rcvAudioSavePath"`
	LogPath           string   `json:"logPath"`
	LogLevel          string   `json:"logLevel"`
	LogPeriod         int      `json:"logPeriod"`
}


func LoadConfigJson(configFileName *string) (ServerConfigJson, error) {
	var parsedResult ServerConfigJson
	file, err := os.Open(*configFileName)
	if err != nil {
		log.Fatal("open error:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&parsedResult)
	if err != nil {
		log.Fatal("json decode error:", err)
	}
	return parsedResult, err
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

// GenerateUnixTrxID Generate UUID v5
func GenerateUnixTrxID() (string, error) {
	uuidCnt = (uuidCnt + 1) % 10000
	uuidString := getMillisTimeFormat(time.Now()) + ":" + localMacAddr + ":" + strconv.Itoa(uuidCnt)
	transactionId := uuid.NewV5(uuid.NamespaceDNS, uuidString)
	return transactionId.String(), nil
}

//Get MAC Address of Server
func getMacAddress() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

// YYYYMMDDhhmmsslll
func getMillisTimeFormat(t time.Time) string {
	// Golang 시간 포멧 기준 2006-01-02 15:04:05, Mon Jan 2 15:04:05 -0700 MST 2006
	timestamp := t.Format("20060102150405")
	return timestamp + strconv.Itoa(t.Nanosecond()/1000000)
}

// YYYYMMDD-hhmmss
func GetFileSaveFormattedTime() string {
	t := time.Now()
	return string(t.Format("20060102-150405"))
}
