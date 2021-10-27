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

	// AUTH error
	HTTP_AUTH_ERR
)