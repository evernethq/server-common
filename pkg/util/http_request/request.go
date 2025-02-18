package http_request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func DoRequest(
	ctx context.Context,
	method, addr string,
	reqBody any,
	opts ...RequestOption,
) (*Response, error) {
	// 默认选项
	options := &requestOptions{
		headers:            make(map[string]string),
		timeout:            3 * time.Second,
		contentType:        ContentTypeJSON,
		allowedStatusCodes: []int{http.StatusOK},
	}

	// 应用选项
	for _, opt := range opts {
		opt(options)
	}

	// 处理JSON请求体
	var bodyReader io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, addr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", string(options.contentType))
	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: options.timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	// 检查状态码是否在允许列表中
	statusAllowed := false
	for _, code := range options.allowedStatusCodes {
		if resp.StatusCode == code {
			statusAllowed = true
			break
		}
	}

	if !statusAllowed {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
