package errors

import (
	"strconv"
)

type CError struct {
	Code    ErrCode
	Message string
}

func NewCError(code ErrCode, msg string) *CError {
	return &CError{
		Code:    code,
		Message: msg,
	}
}

func (e *CError) Error() string {
	return "CODE:" + strconv.Itoa(int(e.Code)) + ", MSG:" + e.Message
}

type ErrCode int

const (
	SUCCESS ErrCode = iota

	// HTTP_SERVER_ERR http.ListenAndServe() error
	HTTP_SERVE_ERR

	// HTTPS_SERVE_ERR https.ListenAndServeTLS() error
	HTTPS_SERVE_ERR

	// Request Header parameter error
	HTTP_REQ_HEADER_PARAM_ERR

	// Request Body parameter error
	HTTP_REQ_PARAMETER_ERR

	// HTTP Invalid method
	HTTP_INVALID_METHOD_ERR

	// HTTP_AUTH_ERR Authorization error
	HTTP_AUTH_ERR

	// HTTP_PREPROCESSING_ERR Cannot processing
	HTTP_PREPROCESSING_ERR

	// HTTP_FLUSHER_ERR Create Flusher http.Flusher error
	HTTP_FLUSHER_ERR

	// HTTP_MULTIPART_PARSE_ERR HTTP Multipart Parse mime.ParseMediaType() error
	HTTP_MULTIPART_PARSE_ERR

	// HTTP_MULTIPART_READ_ERR HTTP Multipart multipartReader.NextPart() error
	HTTP_MULTIPART_READ_ERR

	// HTTP BODY READ NIL
	HTTP_BODY_READ_ERR

	// IOUTIL_READ_ERR ioutil.ReadAll() ioutil.ReadFile() error
	IOUTIL_READ_ERR

	// STT audio format
	STT_FORMAT_ERR

	// STT audio channel format
	STT_CHANNEL_ERR

	// STT audio samplerate
	STT_SAMPLERATE_ERR

	// JSON_MARSHAL_ERR JSON encoding json.Marshal() error
	JSON_MARSHAL_ERR

	// JSON_UNMARSHAL_ERR JSON decoding json.Unmarshal() error
	JSON_UNMARSHAL_ERR

	//Audio Encoding buffer error
	AUDIO_ENCODING_BUFFER_ERR

	//Audio Encoding file error
	AUDIO_ENCODING_FILE_ERR
)
