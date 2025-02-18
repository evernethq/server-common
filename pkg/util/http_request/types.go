package http_request

import (
	"net/http"
	"time"
)

// ContentType 定义请求的内容类型
type ContentType string

const ContentTypeJSON ContentType = "application/json"

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type RequestOption func(*requestOptions)

type requestOptions struct {
	headers            map[string]string
	timeout            time.Duration
	contentType        ContentType
	allowedStatusCodes []int
}

// WithHeaders 设置请求头
func WithHeaders(headers map[string]string) RequestOption {
	return func(o *requestOptions) {
		o.headers = headers
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) RequestOption {
	return func(o *requestOptions) {
		o.timeout = timeout
	}
}

// WithAllowedStatusCodes 配置允许的状态码
func WithAllowedStatusCodes(codes ...int) RequestOption {
	return func(o *requestOptions) {
		o.allowedStatusCodes = codes
	}
}
