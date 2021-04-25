package requests

import (
	"encoding/json"
	"net/http"
	urllib "net/url"
)

type Response struct {
	Request  *request
	History  []*Response
	Encoding string
	Content  []byte
	text     string         // lazy init after calling Text()
	cookies  []*http.Cookie // lazy init after calling Cookies()
	raw      *http.Response
}

// Text 方法将以文本形式读取响应的二进制数据
func (r *Response) Text() string {
	if r.text != "" {
		return r.text
	}
	text := ConvertBytes(&r.Content, "UTF-8")
	r.text = text
	return text
}

// Json 方法将响应转为 JSON
func (r *Response) Json(obj interface{}) error {
	return json.Unmarshal(r.Content, obj)
}

// readUnencodedBodyClose 方法将读取响应的 Body 并返回字节数组指针 *[]byte
func (r *Response) readBodyClose() error {
	var (
		content []byte
		err error
	)
	switch r.raw.Header.Get("Content-Encoding") {
	case "gzip":
		content, err = readGzipBodyClose(r.raw.Body)
	default:
		content, err = readUnencodedBodyClose(r.raw.Body)
	}

	if err != nil {
		return err
	}

	r.Content = content
	return nil
}

// StatusCode: 获取响应状态码
func (r *Response) StatusCode() int {
	return r.raw.StatusCode
}

// URL 获取本次请求的链接
func (r *Response) URL() *urllib.URL {
	return r.raw.Request.URL
}

// URL 获取本次响应返回的 Header
func (r *Response) Headers() http.Header {
	return r.raw.Header
}

// 获取本次 Response 的 cookies
// 包含所有 3xx 跳转的 cookies
func (r *Response) Cookies() []*http.Cookie {
	if r.cookies != nil {
		return r.cookies
	}

	cookies := make([]*http.Cookie, 0)
	for _, his := range r.History {
		cookies = append(cookies, his.raw.Cookies()...)
	}
	cookies = append(cookies, r.raw.Cookies()...)

	r.cookies = cookies
	return cookies
}

func (r *Response) Raw() *http.Response {
	return r.raw
}
