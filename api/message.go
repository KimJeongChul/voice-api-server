package api

type MsgSpeechMetaData struct {
	Encoding       string               `json:"encoding"`
	TargetLanguage string               `json:"targetLanguage"`
	SttMode        int                  `json:"sttMode"`
	EncodingOpt    MsgSpeechEncodingOpt `json:"encodingOpt"`
}

// MsgSpeechEncodingOpt
// - samplefmt: s16le, f32le
// - channel: 1(mono), 2(stereo)
// - samplerate: 16000, 44100, 48000
type MsgSpeechEncodingOpt struct {
	Channel    int    `json:"channel"`
	SampleRate int    `json:"sampleRate"`
	SampleFmt  string `json:"sampleFmt"`
}
