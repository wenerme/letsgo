package bg

import (
	"fmt"
	"strings"
)

type Created struct {
	ID uint `json:"id"`
}

type General struct {
	Code    int `json:"code"`
	Message string `json:"message,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Error struct {
	// HTML 状态码
	Status  int `json:"status,omitempty"`
	// 错误码,如果没有则为 Code
	Code    int `json:"code,omitempty"`
	// 错误消息
	Message string `json:"message,omitempty"`
	// 请求路径
	Path    string `json:"path,omitempty"`
	cause   error
}

func (self Error)Error() string {
	return self.Message
}
func (self Error)Cause() error {
	return self.cause
}

func CreateError(err interface{}) (e Error) {
	switch v := err.(type){
	case Error:
		e = v
	case error:
		e.Message = fmt.Sprint(err)
		e.cause = v
	case string:
		e.Message = v
	default:
		e.Message = fmt.Sprint(err)
	}
	e.Message = strings.TrimSpace(e.Message)
	e.Status = 500
	e.Code = e.Status
	return e
}
func CreateErrorWithStatus(status int, err interface{}) (e Error) {
	e = CreateError(err)
	e.Status = status
	e.Code = status
	return
}
func CreateErrorBadRequest(err interface{}) Error {
	return CreateErrorWithStatus(400, err)
}
