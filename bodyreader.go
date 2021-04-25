package requests

import (
	"compress/gzip"
	"io"
	"io/ioutil"
)

type BodyReaderFunc func(body io.ReadCloser) ([]byte, error)


func readUnencodedBodyClose(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func readGzipBodyClose(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	reader, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return content, nil
}





