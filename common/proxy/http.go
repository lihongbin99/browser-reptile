package proxy

import (
	"browser-reptile/common/utils"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type HttpHandle interface {
	NeedParse(host string, port uint16) bool
	RequestNeedParseHeader(proxy *HttpProxy) bool
	RequestParseHeader(proxy *HttpProxy)
	RequestNeedParseBody(proxy *HttpProxy) bool
	RequestParseBody(proxy *HttpProxy, body []byte) []byte
	ResponseNeedParseHeader(proxy *HttpProxy) bool
	ResponseParseHeader(proxy *HttpProxy)
	ResponseNeedParseBody(proxy *HttpProxy) bool
	ResponseParseBody(proxy *HttpProxy, body []byte) []byte
}

type HttpProxy struct {
	*Proxy
	RequestMethod   string
	RequestPath     string
	RequestProtocol string
	RequestHeader   map[string][]string

	ResponseProtocol    string
	ResponseCode        string
	ResponseCodeMessage string
	ResponseHeader      map[string][]string

	Handles  []HttpHandle
	waitChan chan uint8
}

func NewHttpProxy(socksProxy *Proxy, handles []HttpHandle, waitChan chan uint8) *HttpProxy {
	return &HttpProxy{
		socksProxy,
		"", "", "", nil,
		"", "", "", nil,
		handles, waitChan,
	}
}

func (that *HttpProxy) HttpRequest() {
	defer func() {
		that.waitChan <- 0
	}()

	// 获取请求协议
	protocolErr := true
	protocolErr, that.RequestMethod, that.RequestPath, that.RequestProtocol = utils.HttpReadProtocol(that.ClientConn)
	if protocolErr {
		return
	}
	// TODO 处理因为 Transfer-Encoding 引发的问题
	that.RequestProtocol = "HTTP/1.0"

	// 发送请求协议
	_, _ = that.ServerConn.Write([]byte(fmt.Sprintf("%s %s %s\r\n", that.RequestMethod, that.RequestPath, that.RequestProtocol)))

	// 获取请求头
	requestHeader := utils.HttpReadHeader(that.ClientConn)
	if requestHeader == nil {
		return
	}
	that.RequestHeader = requestHeader

	// 处理请求头
	for _, handle := range that.Handles {
		if handle.RequestNeedParseHeader(that) {
			handle.RequestParseHeader(that)
		}
	}

	// 发送请求头
	headerStr := ""
	for headerName, headerValues := range that.RequestHeader {
		for _, headerValue := range headerValues {
			headerStr += headerName + ": " + headerValue + "\r\n"
		}
	}
	headerStr += "\r\n"
	_, _ = that.ServerConn.Write([]byte(headerStr))

	// 处理请求体
	var body []byte = nil
	for _, handle := range that.Handles {
		if handle.RequestNeedParseBody(that) {
			if body == nil {
				var err error = nil
				body, err = ioutil.ReadAll(that.ClientConn)
				if err != nil {
					return
				}
			}
			body = handle.RequestParseBody(that, body)
		}
	}

	// 发送请求体
	if body != nil {
		_, _ = that.ServerConn.Write(body)
	} else {
		buf := make([]byte, 64*1024)
		for {
			readLen, err := that.ClientConn.Read(buf)
			if (err != nil && err != io.EOF) || readLen <= 0 {
				return
			}
			if _, err = that.ServerConn.Write(buf[:readLen]); err != nil {
				return
			}
			if err == io.EOF {
				break
			}
		}
	}
}

func (that *HttpProxy) HttpResponse() {
	defer func() {
		_ = that.ServerConn.Close()
		_ = that.ClientConn.Close()
		that.waitChan <- 0
	}()

	// 获取响应协议
	protocolErr := true
	protocolErr, that.ResponseProtocol, that.ResponseCode, that.ResponseCodeMessage = utils.HttpReadProtocol(that.ServerConn)
	if protocolErr {
		return
	}

	// 发送响应协议
	_, _ = that.ClientConn.Write([]byte(fmt.Sprintf("%s %s %s\r\n", that.ResponseProtocol, that.ResponseCode, that.ResponseCodeMessage)))

	// 获取响应头
	responseHeader := utils.HttpReadHeader(that.ServerConn)
	if responseHeader == nil {
		return
	}
	that.ResponseHeader = responseHeader

	// 处理响应头
	for _, handle := range that.Handles {
		if handle.ResponseNeedParseHeader(that) {
			handle.ResponseParseHeader(that)
		}
	}

	// 发送响应头
	headerStr := ""
	for headerName, headerValues := range that.ResponseHeader {
		for _, headerValue := range headerValues {
			headerStr += headerName + ": " + headerValue + "\r\n"
		}
	}
	headerStr += "\r\n"
	_, _ = that.ClientConn.Write([]byte(headerStr))

	// 处理响应体
	gzip := false
	if values, exist := that.ResponseHeader["Content-Encoding"]; exist {
		for _, value := range values {
			if strings.Contains(value, "gzip") {
				gzip = true
				break
			}
		}
	}
	var body []byte = nil
	for _, handle := range that.Handles {
		if handle.ResponseNeedParseBody(that) {
			if body == nil {
				var err error = nil
				body, err = ioutil.ReadAll(that.ServerConn)
				if err != nil {
					return
				}
				if gzip {
					// 解压 gzip
					if body, err = utils.Gunzip(body); err != nil {
						return
					}
				}
			}
			body = handle.ResponseParseBody(that, body)
		}
	}

	// 发送响应体
	if body != nil {
		if gzip {
			// 压缩 gzip
			var err error = nil
			body, err = utils.Gzip(body)
			if err != nil {
				return
			}
		}
		_, _ = that.ClientConn.Write(body)
	} else {
		buf := make([]byte, 64*1024)
		for {
			readLen, err := that.ServerConn.Read(buf)
			if (err != nil && err != io.EOF) || readLen <= 0 {
				return
			}
			if _, err = that.ClientConn.Write(buf[:readLen]); err != nil {
				return
			}
			if err == io.EOF {
				break
			}
		}
	}
}
