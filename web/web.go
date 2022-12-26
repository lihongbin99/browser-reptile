package web

import (
	"browser-reptile/common/config"
	"browser-reptile/common/utils"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func Main(conn net.Conn) {

	// 获取请求协议
	protocolErr, method, path, protocol := utils.HttpReadProtocol(conn)
	if protocolErr {
		return
	}

	headers, headerSrcs := utils.HttpReadHeader(conn)
	if headers == nil {
		return
	}

	data := make([]byte, 0)
	var err error = nil
	if method != "GET" {
		data, err = ioutil.ReadAll(conn)
		if err != nil {
			return
		}
	}

	doMain(conn, method, path, protocol, headers, headerSrcs, data)
}

var (
	cerBuf []byte = nil
)

func doMain(conn net.Conn, method, path, protocol string, headers map[string][]string, headerSrcs map[string]string, data []byte) {
	if path == "/proxy.cer" {
		// 获取数据
		if cerBuf == nil {
			var err error = nil
			cerBuf, err = os.ReadFile(config.TlsCertFile)
			if err != nil {
				return
			}
		}

		// 发送响应协议
		response := make([]byte, 0)
		response = append(response, protocol...)
		response = append(response, " 200 OK\r\n"...)
		response = append(response, fmt.Sprintf("Content-Length: %d\r\n", len(cerBuf))...)
		response = append(response, "Connection: close\r\n"...)
		response = append(response, "Cache-Control: max-age=0\r\n"...)
		response = append(response, "Content-Type: application/x-x509-ca-cert\r\n"...)
		response = append(response, "\r\n"...)
		response = append(response, cerBuf...)
		response = append(response, "\n"...)

		// 发送响应体
		_, _ = conn.Write(response)
	}
}
