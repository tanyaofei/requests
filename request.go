package requests

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	urllib "net/url"
	"strings"
)

type request struct {
	// HTTP 请求方法
	Method string

	// HTTP 请求链接
	URL *urllib.URL

	// HTTP 请求头
	Headers Headers

	// HTTP multipart/form-data or application/x-www-form-urlencoded
	// 如果 Json 不为 nil, Data 将会失效
	// 如果 Files 不为 nil, 则请求头将会变为 multipart/form-data
	// 如果 Files 为 nil, 请求头将会为 application/x-www-form-urlencoded
	Data Data

	// Http multipart/form-data file
	// 如果 Json != nil, Files 将会失效
	Files Files

	// HTTP 链接参数
	Query Query

	// Cookies
	Cookies []*http.Cookie

	// HTTP application/json
	Json Json

	// 原始的请求
	raw *http.Request
}

// prepare 用与来构建内置库结构体(struct) http.Request
// 如果 request.Json 不为空, 则会将 "Content-Type": "application/json" 写入 http.Request 的 Header(http.Header) 属性
// 并且 将 request.Json 写入到 http.Request 的 Body
// 如果 request.Data 不为空, 则会将 "Content-Type": "multipart/form-data; boundary=xx" 写入 http.Request 的 Header(http.Header) 属性
// 并且 将 request.Data 写入到 http.Request 的 Body 中
func (r *request) prepare() error {
	if err := r.prepareRawAndBody(); err != nil {
		return errors.WithStack(err)
	}

	r.prepareHeaders()
	r.prepareCookies()
	r.prepareParams()
	// set back
	r.URL = r.raw.URL

	return nil
}

// prepareRawAndBody 应该要第一个调用, 调用后才会创建 raw
func (r *request) prepareRawAndBody() error {
	var (
		body        io.Reader
		contentType = r.Headers["Content-Type"]
	)

	// body
	if r.Json != nil {
		// application/json body
		jsonBytes, err := json.Marshal(r.Json)
		if err != nil {
			return errors.WithStack(err)
		}
		body = bytes.NewBuffer(jsonBytes)
		if contentType == "" {
			contentType = "application/json"
		}

	} else if r.Files != nil {
		// 有 Files
		// Content-Type: multipart/form-data
		buffer := new(bytes.Buffer)
		writer := multipart.NewWriter(buffer)
		for k, v := range r.Files {
			err := writer.WriteField(k, string(*v))
			if err != nil {
				return errors.WithStack(err)
			}
		}

		// 有 Files 也有 Data
		if r.Data != nil {
			for k, v := range r.Data {
				err := writer.WriteField(k, v)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
		_ = writer.Close()
		body = io.Reader(buffer)
		contentType = writer.FormDataContentType()

	} else if r.Files == nil && r.Data != nil {
		// 发送 data 但是不发送文件
		// Content-Type: application/x-www-form-urlencoded
		data := urllib.Values{}
		for k, v := range r.Data {
			data.Set(k, v)
		}
		body = strings.NewReader(data.Encode())
		if contentType == "" {
			contentType = "application/x-www-form-urlencoded"
		}
	}

	raw, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return err
	}

	if contentType != "" {
		raw.Header.Set("Content-Type", contentType)
		if r.Headers == nil {
			r.Headers = make(Headers)
		}
		r.Headers["Content-Type"] = contentType
	}

	if r.Query != nil {
		q := r.URL.Query()
		for k, v := range r.Query {
			q.Add(k, v)
		}
		raw.URL.RawQuery = q.Encode()
	}

	r.raw = raw
	return nil
}

// 在调用 prepareHeaders 应先创建 raw
func (r *request) prepareHeaders() {
	if r.Headers != nil {
		for k, v := range r.Headers {
			r.raw.Header.Set(k, v)
		}
	}
}

// 在调用 prepareParams 应先创建 raw
func (r *request) prepareParams() {
	if r.Query != nil {
		for k, v := range r.Query {
			r.raw.URL.Query().Set(k, v)
		}
	}
}

// 在调用 prepareCookies 应先创建 raw
func (r *request) prepareCookies() {
	if r.Cookies != nil {
		for _, cookie := range r.Cookies {
			if cookie.Domain != r.URL.Host {
				continue
			}
			if cookie.Path != "" && !strings.HasPrefix(r.URL.Path, cookie.Path) {
				continue
			}
			r.raw.AddCookie(cookie)
		}
	}
}
