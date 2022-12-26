package utils

import (
	"fmt"
	"io"
	"net"
)

type GzipChunkedReader struct {
	skip int
	data []byte
}

func NewGzipChunkedReader(data []byte) *GzipChunkedReader {
	return &GzipChunkedReader{0, data}
}

func (that *GzipChunkedReader) Read(buf []byte) (n int, err error) {
	length := int(buf[that.skip])
	if length == 0 {
		return -1, io.EOF
	}
	fmt.Println(length)

	that.skip++
	for i := 0; i < length; i++ {
		buf[i] = that.data[that.skip+i]
	}
	that.skip += length

	return length, nil
}

func ReadN(conn net.Conn, buf []byte, len int) error {
	readLen := 0
	for len-readLen > 0 {
		n, err := conn.Read(buf[readLen:len])
		if err != nil {
			return err
		}
		readLen += n
	}
	return nil
}
