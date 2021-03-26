package requests

const (
	version = "0.1"
)

func Get(url string, options ...interface{}) (*Response, error) {
	return Request("GET", url, options...)
}

func Post(url string, options ...interface{}) (*Response, error) {
	return Request("POST", url, options...)
}

func Put(url string, options ...interface{}) (*Response, error) {
	return Request("PUT", url, options...)
}

func Delete(url string, options ...interface{}) (*Response, error) {
	return Request("DELETE", url, options...)
}

func Patch(url string, options ...interface{}) (*Response, error) {
	return Request("PATCH", url, options...)
}

func Head(url string, options ...interface{}) (*Response, error) {
	return Request("HEAD", url, options...)
}

func Options(url string, options ...interface{}) (*Response, error) {
	return Request("OPTIONS", url, options...)
}

func Request(method, url string, options ...interface{}) (*Response, error) {
	session, err := NewSession(options...)
	if err != nil {
		return nil, err
	}
	return session.Request(method, url, options...)
}
