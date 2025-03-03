package errors

import krErr "github.com/go-kratos/kratos/v2/errors"

// TimeOut 超时
var TimeOut = krErr.New(504, "TIMEOUT", "")

// Unknown 未知错误
var Unknown = func(msg string) error {
	return krErr.New(500, "UNKNOWN", msg)
}
