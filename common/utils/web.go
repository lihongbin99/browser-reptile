package utils

import (
	"io"
	"net"
	"strings"
)

func HttpReadProtocol(conn net.Conn) (err bool, s1, s2, s3 string) {
	buf := make([]byte, 1)
	data := make([]byte, 0, 64*1024)
	for {
		readLen, e := conn.Read(buf[0:1])
		if (e != nil && e != io.EOF) || readLen <= 0 {
			return true, "", "", ""
		}
		data = append(data, buf[0])
		if len(data) >= 2 {
			if string(data[len(data)-2:]) == "\r\n" {
				break
			}
		}
		if e == io.EOF {
			return true, "", "", ""
		}
	}

	httpProtocols := strings.Split(string(data[:len(data)-2]), " ")
	if len(httpProtocols) >= 3 {
		return false, httpProtocols[0], httpProtocols[1], httpProtocols[2]
	}

	return true, "", "", ""
}

func HttpReadHeader(conn net.Conn) map[string][]string {
	buf := make([]byte, 1)
	data := make([]byte, 0, 64*1024)
	for {
		readLen, err := conn.Read(buf[0:1])
		if (err != nil && err != io.EOF) || readLen <= 0 {
			return nil
		}
		data = append(data, buf[0])
		if len(data) >= 4 {
			if string(data[len(data)-4:]) == "\r\n\r\n" {
				break
			}
		}
		if err == io.EOF {
			return nil
		}
	}

	// 解析请求头
	headers := make(map[string][]string)
	headerLine := strings.Split(string(data), "\r\n")
	for _, str := range headerLine {
		if str != "" {
			// 处理请求头
			headerIndex := strings.Index(str, ": ")
			if headerIndex >= 0 {
				headerName := str[:headerIndex]
				headerValue := str[headerIndex+2:]
				if values, ok := headers[headerName]; ok {
					headers[headerName] = append(values, headerValue)
				} else {
					values = make([]string, 1)
					values[0] = headerValue
					headers[headerName] = values
				}
			}
		}
	}
	return headers
}
