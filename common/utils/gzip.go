package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Gunzip(data []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(gzipReader)
}

func Gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gzipWriter := gzip.NewWriter(&b)
	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}
	_ = gzipWriter.Flush()
	_ = gzipWriter.Close()
	return b.Bytes(), nil
}
