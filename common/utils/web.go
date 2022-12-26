package utils

import (
	"io"
	"net"
	"strconv"
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

func HttpReadHeader(conn net.Conn) (map[string][]string, map[string]string) {
	buf := make([]byte, 1)
	data := make([]byte, 0, 64*1024)
	for {
		readLen, err := conn.Read(buf[0:1])
		if (err != nil && err != io.EOF) || readLen <= 0 {
			return nil, nil
		}
		data = append(data, buf[0])
		if len(data) >= 4 {
			if string(data[len(data)-4:]) == "\r\n\r\n" {
				break
			}
		}
		if err == io.EOF {
			return nil, nil
		}
	}

	// 解析请求头
	headers := make(map[string][]string)
	headerSrcs := make(map[string]string)
	headerLine := strings.Split(string(data), "\r\n")
	for _, str := range headerLine {
		if str != "" {
			// 处理请求头
			headerIndex := strings.Index(str, ": ")
			if headerIndex >= 0 {
				headerName := str[:headerIndex]
				headerValue := str[headerIndex+2:]

				headerLowName := strings.ToLower(headerName)
				headerSrcs[headerLowName] = headerName

				if values, ok := headers[headerLowName]; ok {
					headers[headerLowName] = append(values, headerValue)
				} else {
					values = make([]string, 1)
					values[0] = headerValue
					headers[headerLowName] = values
				}
			}
		}
	}
	return headers, headerSrcs
}

func HttpChunked(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1)
	data := make([]byte, 0, 64)
	result := make([]byte, 0)

	for {
		readLen, err := conn.Read(buf[0:1])
		if err != nil || readLen <= 0 {
			return nil, err
		}
		// 读取长度
		data = append(data, buf[0])
		if len(data) >= 3 {
			if string(data[len(data)-2:]) == "\r\n" {
				// 解析长度
				contentLen64, err := strconv.ParseInt(string(data[0:len(data)-2]), 16, 32)
				if err != nil {
					return nil, err
				}
				contentLen := int(contentLen64)

				if contentLen == 0 {
					_ = ReadN(conn, data, 2) // 最后的\r\n
					break
				}

				// 读取响应体
				temp := make([]byte, contentLen+2) // +2是因为内容后面带了\r\n
				if err = ReadN(conn, temp, contentLen+2); err != nil {
					return nil, err
				}
				result = append(result, temp[0:contentLen]...)
				data = make([]byte, 0, 64)
			}
		}
	}

	return result, nil
}
