package bg

import (
	"fmt"
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
}

func (self Error)Error() string {
	return self.Message
}

func CreateError(err interface{}) (e Error) {
	if v, ok := err.(error); ok {
		e.Message = v.Error()
	} else {
		e.Message = fmt.Sprint(err)
	}
	return e
}
