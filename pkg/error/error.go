package errno

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var _ Error = (*err)(nil)

type Error interface {
	p()
	// WithErr 设置错误信息
	WithErr(err error) Error
	// GetBusinessCode 获取 Business Code
	GetBusinessCode() int
	// GetHttpCode 获取 HTTP Code
	GetHttpCode() int
	// GetMsg 获取 Msg
	GetMsg() string
	// ToString 返回 JSON 格式的错误详情
	ToString() string
}

type err struct {
	HttpCode     int
	BusinessCode int
	Message      string
	Err          error
}

func NewError(httpCode, businessCode int, msg string) Error {
	return &err{
		HttpCode:     httpCode,
		BusinessCode: businessCode,
		Message:      msg,
	}
}

func (e *err) p() {}

func (e *err) WithErr(err error) Error {
	e.Err = errors.WithStack(err)
	return e
}

func (e *err) GetHttpCode() int {
	return e.HttpCode
}

func (e *err) GetBusinessCode() int {
	return e.BusinessCode
}

func (e *err) GetMsg() string {
	return e.Message
}

// ToString 返回 JSON 格式的错误详情
func (e *err) ToString() string {
	err := &struct {
		HttpCode     int    `json:"code"`
		BusinessCode int    `json:"error_code"`
		Message      string `json:"message"`
	}{
		HttpCode:     e.HttpCode,
		BusinessCode: e.BusinessCode,
		Message:      e.Message,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}
