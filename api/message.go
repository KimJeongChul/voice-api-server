package api

type MsgSpeechMetaData struct {
	Encoding       string               `json:"encoding"`
	TargetLanguage string               `json:"targetLanguage"`
	SttMode        int                  `json:"sttMode"`
	EncodingOpt    MsgSpeechEncodingOpt `json:"encodingOpt"`
}

// MsgSpeechEncodingOpt
// - samplefmt: S16LE, F32LE
// - channel: 1(mono), 2(stereo)
// - samplerate: 16000, 44100, 48000
type MsgSpeechEncodingOpt struct {
	Channel    int    `json:"channel"`
	SampleRate int    `json:"sampleRate"`
	SampleFmt  string `json:"sampleFmt"`
}

type FileSpeechRecognizeRes struct {
	TransactionId string						    `json:"transactionID"`
	EventTime     string							`json:"eventTime"`
	SttStatus  	  string                            `json:"sttStatus"`
	SttResults    *[]MsgSpeechRecognizeResSttResult `json:"sttResults,omitempty"`
	SttInfo       *MsgSpeechRecognizeResSttInfo     `json:"sttInfo,omitempty"`
}

// MsgVoiceRecognizeRes Response Body
type MsgSpeechRecognizeRes struct {
	ResultType string                          `json:"resultType"`
	SttResult  *MsgSpeechRecognizeResSttResult `json:"speechResult,omitempty"`
	SttInfo    *MsgSpeechRecognizeResSttInfo   `json:"speechInfo,omitempty"`
	ErrCode    string                          `json:"errCode,omitempty"`
}

type MsgSpeechRecognizeResSttResult struct {
	Text      string  `json:"text"`
	StartTime float32 `json:"startTime,omitempty"`
	EndTime   float32 `json:"endTime,omitempty"`
}

type MsgSpeechRecognizeResSttInfo struct {
	ReqFileSize  int    `json:"reqFileSize"`
	TransCodec   string `json:"transCodec"`
	ConvFileSize int    `json:"convFileSize"`
	SttInputTime int    `json:"speechInputTime"`
}
