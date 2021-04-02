package requests

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net"
	"net/http"
	urllib "net/url"
	"strings"
	"time"
)

type session struct {
	Headers  Headers
	Params   Query
	Verify   Verify
	Redirect Redirect
	Cookies  []*http.Cookie
	Proxy    Proxy
	raw      *http.Client
}

func (s *session) Get(url string, options ...interface{}) (*Response, error) {
	return s.Request("GET", url, options...)
}

func (s *session) Post(url string, options ...interface{}) (*Response, error) {
	return s.Request("POST", url, options...)
}

func (s *session) Put(url string, options ...interface{}) (*Response, error) {
	return s.Request("PUT", url, options...)
}

func (s *session) Delete(url string, options ...interface{}) (*Response, error) {
	return s.Request("DELETE", url, options...)
}

func (s *session) Head(url string, options ...interface{}) (*Response, error) {
	return s.Request("HEAD", url, options...)
}

func (s *session) Patch(url string, options ...interface{}) (*Response, error) {
	return s.Request("PATCH", url, options...)
}

func (s *session) Options(url string, options ...interface{}) (*Response, error) {
	return s.Request("OPTIONS", url, options...)
}

// Request 方法将创建一个 request, 并将 method, url, options 写入该 request 中
func (s *session) Request(method, url string, options ...interface{}) (*Response, error) {
	var (
		req      = new(request)
		err      error
		redirect = s.Redirect
	)
	req.URL, err = urllib.Parse(url)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to parse URL: %s", url)
	}

	req.Method = method
	req.Headers = s.Headers
	req.Query = s.Params
	req.Cookies = s.Cookies

	for _, opt := range options {
		if opt == nil {
			continue
		}
		switch opt.(type) {
		case Data:
			req.Data = opt.(Data)
		case Json:
			req.Json = opt.(Json)
		case Headers:
			if req.Headers == nil {
				req.Headers = opt.(Headers)
			} else {
				for k, v := range opt.(Headers) {
					req.Headers[k] = v
				}
			}
		case http.Header:
			if req.Headers == nil {
				s.Headers = make(Headers)
			}
			for k, v := range opt.(http.Header) {
				s.Headers[k] = v[0]
			}
		case Query:
			if req.Query == nil {
				req.Query = opt.(Query)
			} else {
				for k, v := range opt.(Query) {
					req.Query[k] = v
				}
			}
		case Cookies:
			cookies := opt.(Cookies)
			if req.Cookies == nil {
				req.Cookies = make([]*http.Cookie, len(cookies))
			}
			for k, v := range cookies {
				req.Cookies = append(req.Cookies, &http.Cookie{Name: k, Value: v})
			}
		case []*http.Cookie:
			if req.Cookies == nil {
				req.Cookies = opt.([]*http.Cookie)
			} else {
				req.Cookies = append(req.Cookies, opt.([]*http.Cookie)...)
			}
		case Redirect:
			redirect = opt.(Redirect)
		}

	}

	// sendAndSetCookies http with inner http lib
	resp, err := s.doRequest(req, nil, redirect)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

// sendAndSetCookies 将会发送 http.Request 请求
// 并将该请求的响应的 Set-Cookie 写入到 session
func (s *session) sendAndSetCookies(rawReq *http.Request) (*http.Response, error) {
	rawResp, err := s.raw.Do(rawReq)
	if err != nil {
		return nil, err
	}

	// set cookie into session from Response
	s.Cookies = append(s.Cookies, rawResp.Cookies()...)
	return rawResp, nil
}

// doRequest 方法用于调用 sendAndSetCookies 方法发送请求, 同时将会处理 3xx 跳转
// 收集所有跳转记录写入 Response 并最后返回 Response
func (s *session) doRequest(req *request, lastResp *Response, redirect Redirect) (*Response, error) {
	err := req.prepare()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rawResp, err := s.sendAndSetCookies(req.raw)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	currentResp := &Response{Request: req, raw: rawResp}
	// 收集 history
	if lastResp != nil {
		if lastResp.History != nil {
			currentResp.History = append(lastResp.History[:], lastResp)
		} else {
			currentResp.History = []*Response{lastResp}
		}
	}

	// 非 3xx 跳转响应
	if rawResp.Header.Get("Location") == "" || !redirect {
		return currentResp, nil
	}

	// 3xx 跳转响应
	location := rawResp.Header.Get("Location")
	if !strings.HasPrefix(location, "http") {
		location = req.URL.Scheme + "://" + req.URL.Host + location
	}

	nextURL, err := urllib.Parse(location)
	if err != nil {
		return nil, err
	}
	nextHeader := Headers{}
	if s.Headers != nil {
		for k, v := range nextHeader {
			nextHeader[k] = v
		}
	}
	nextHeader["Referer"] = req.URL.String()
	nextReq := &request{
		URL:     nextURL,
		Headers: nextHeader,
		Cookies: s.Cookies,
		Query:   s.Params,
	}

	return s.doRequest(nextReq, currentResp, redirect)
}

func NewSession(options ...interface{}) (*session, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	s := &session{
		Redirect: true,
		Verify:   true,
		raw: &http.Client{
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}},
	}

	for _, opt := range options {
		if opt == nil {
			continue
		}
		switch opt.(type) {
		case Headers:
			s.Headers = opt.(Headers)
		case http.Header:
			s.Headers = make(Headers)
			for k, v := range opt.(http.Header) {
				s.Headers[k] = v[0]
			}
		case Query:
			s.Params = opt.(Query)
		case []*http.Cookie:
			s.Cookies = opt.([]*http.Cookie)
		case Cookies:
			s.Cookies = make([]*http.Cookie, 0)
			for k, v := range opt.(Cookies) {
				s.Cookies = append(s.Cookies, &http.Cookie{Name: k, Value: v})
			}
		case Verify:
			if !opt.(Verify) {
				transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		case Proxy:
			s.Proxy = opt.(Proxy)
			if string(s.Proxy) != "" {
				transport.Proxy = func(r *http.Request) (*urllib.URL, error) {
					return urllib.Parse(string(s.Proxy))
				}
			}
		case Redirect:
			s.Redirect = opt.(Redirect)
		}

	}

	return s, nil
}
