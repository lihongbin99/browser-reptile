package proxy

import (
	"browser-reptile/common/utils"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
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
	RequestMethod    string
	RequestPath      string
	RequestProtocol  string
	RequestHeader    map[string][]string
	RequestSrcHeader map[string]string

	ResponseProtocol    string
	ResponseCode        string
	ResponseCodeMessage string
	ResponseHeader      map[string][]string
	ResponseSrcHeader   map[string]string

	Handles  []HttpHandle
	waitChan chan uint8
}

func NewHttpProxy(socksProxy *Proxy, handles []HttpHandle, waitChan chan uint8) *HttpProxy {
	return &HttpProxy{
		socksProxy,
		"", "", "", nil, nil,
		"", "", "", nil, nil,
		handles, waitChan,
	}
}

func (that *HttpProxy) HttpRequest() {
	defer func() {
		that.waitChan <- 0
	}()

	// 循环处理 HTTP1.1 的长连接
	for {
		// 获取请求协议
		protocolErr := true
		protocolErr, that.RequestMethod, that.RequestPath, that.RequestProtocol = utils.HttpReadProtocol(that.ClientConn)
		if protocolErr {
			_ = that.ClientConn.Close()
			_ = that.ServerConn.Close()
			return
		}

		// 发送请求协议
		_, _ = that.ServerConn.Write([]byte(fmt.Sprintf("%s %s %s\r\n", that.RequestMethod, that.RequestPath, that.RequestProtocol)))

		// 获取请求头
		requestHeader, requestSrcHeader := utils.HttpReadHeader(that.ClientConn)
		if requestHeader == nil {
			_ = that.ClientConn.Close()
			_ = that.ServerConn.Close()
			return
		}
		that.RequestHeader = requestHeader
		that.RequestSrcHeader = requestSrcHeader

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
				headerStr += that.RequestSrcHeader[headerName] + ": " + headerValue + "\r\n"
			}
		}
		headerStr += "\r\n"
		_, _ = that.ServerConn.Write([]byte(headerStr))

		// 处理 HTTP1.1 的长连接情况
		keepAliveAndContentLength := false
		contentLength := 0
		keepAliveAndTransferEncoding := ""
		if that.RequestProtocol == "HTTP/1.1" {
			if contentLengthHeader, ok := that.RequestHeader["content-length"]; ok {
				var err error = nil
				if contentLength, err = strconv.Atoi(contentLengthHeader[0]); err == nil {
					keepAliveAndContentLength = true
				}
			}
			if contentLengthHeader, ok := that.RequestHeader["transfer-encoding"]; ok {
				keepAliveAndTransferEncoding = contentLengthHeader[0]
			}
		}

		if that.RequestMethod == "GET" || that.RequestMethod == "OPTIONS" {
			continue
		}
		if keepAliveAndTransferEncoding != "" && !keepAliveAndContentLength {
			fmt.Println("在", that.Host, that.RequestPath, "中出现了Transfer-Encoding:", keepAliveAndTransferEncoding)
		}

		// 获取请求体
		var body []byte = nil
		if keepAliveAndContentLength {
			body = make([]byte, contentLength)
			if err := utils.ReadN(that.ClientConn, body, contentLength); err != nil {
				_ = that.ClientConn.Close()
				_ = that.ServerConn.Close()
				return
			}
		}

		// 处理请求体
		for _, handle := range that.Handles {
			if handle.RequestNeedParseBody(that) {
				if body == nil {
					var err error = nil
					body, err = ioutil.ReadAll(that.ClientConn)
					if err != nil {
						_ = that.ClientConn.Close()
						_ = that.ServerConn.Close()
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
					return
				}
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

	// 循环处理 HTTP1.1 的长连接
	for {
		// 获取响应协议
		protocolErr := true
		protocolErr, that.ResponseProtocol, that.ResponseCode, that.ResponseCodeMessage = utils.HttpReadProtocol(that.ServerConn)
		if protocolErr {
			return
		}

		// 发送响应协议
		_, _ = that.ClientConn.Write([]byte(fmt.Sprintf("%s %s %s\r\n", that.ResponseProtocol, that.ResponseCode, that.ResponseCodeMessage)))

		// 获取响应头
		responseHeader, responseSrcHeader := utils.HttpReadHeader(that.ServerConn)
		if responseHeader == nil {
			return
		}
		that.ResponseHeader = responseHeader
		that.ResponseSrcHeader = responseSrcHeader

		// 处理响应头
		for _, handle := range that.Handles {
			if handle.ResponseNeedParseHeader(that) {
				handle.ResponseParseHeader(that)
			}
		}

		/*暂时不发送响应头, 因为如果修改了响应体的话可能需要修改响应头*/

		// 获取响应体参数
		gzip := false
		if values, exist := that.ResponseHeader["content-encoding"]; exist {
			for _, value := range values {
				if strings.Contains(value, "gzip") {
					gzip = true
					break
				}
			}
		}

		// 获取响应体参数
		keepAliveAndContentLength := false
		contentLength := 0
		keepAliveAndTransferEncoding := ""
		if that.RequestProtocol == "HTTP/1.1" {
			if contentLengthHeader, ok := that.ResponseHeader["content-length"]; ok {
				var err error = nil
				if contentLength, err = strconv.Atoi(contentLengthHeader[0]); err == nil {
					keepAliveAndContentLength = true
				}
			}
			if contentLengthHeader, ok := that.ResponseHeader["transfer-encoding"]; ok {
				keepAliveAndTransferEncoding = contentLengthHeader[0]
			}
		}

		// 获取响应体
		var body []byte = nil
		if keepAliveAndContentLength {
			body = make([]byte, contentLength)
			if err := utils.ReadN(that.ServerConn, body, contentLength); err != nil {
				return
			}
		} else if keepAliveAndTransferEncoding != "" {
			if keepAliveAndTransferEncoding == "chunked" {
				var err error = nil
				if body, err = utils.HttpChunked(that.ServerConn); err != nil {
					return
				}
			} else if keepAliveAndTransferEncoding == "compress" {
				fmt.Println("访问", that.Host, that.RequestPath, "时出现了Transfer-Encoding:", keepAliveAndTransferEncoding)
				return
			} else if keepAliveAndTransferEncoding == "deflate" {
				fmt.Println("访问", that.Host, that.RequestPath, "时出现了Transfer-Encoding:", keepAliveAndTransferEncoding)
				return
			} else if keepAliveAndTransferEncoding == "gzip" {
				fmt.Println("访问", that.Host, that.RequestPath, "时出现了Transfer-Encoding:", keepAliveAndTransferEncoding)
				return
			} else {
				fmt.Println("访问", that.Host, that.RequestPath, "时出现了未知的Transfer-Encoding:", keepAliveAndTransferEncoding)
				return
			}
		}

		if body != nil {
			if gzip {
				// 解压 gzip
				var err error = nil
				if body, err = utils.Gunzip(body); err != nil {
					return
				}
			}
		}

		// 处理响应体
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

		if body != nil {
			if gzip {
				// 压缩 gzip
				var err error = nil
				body, err = utils.Gzip(body)
				if err != nil {
					return
				}
			}
			// 修改 Content-Length
			if values, exist := that.ResponseHeader["content-length"]; exist {
				values = make([]string, 1)
				values[0] = strconv.Itoa(len(body))
				that.ResponseHeader["content-length"] = values
			}
		}

		// 发送响应头
		headerStr := ""
		for headerName, headerValues := range that.ResponseHeader {
			for _, headerValue := range headerValues {
				headerStr += that.ResponseSrcHeader[headerName] + ": " + headerValue + "\r\n"
			}
		}
		headerStr += "\r\n"
		_, _ = that.ClientConn.Write([]byte(headerStr))

		// 发送响应体
		if body != nil {
			if keepAliveAndTransferEncoding == "chunked" {
				len64 := strconv.FormatInt(int64(len(body)), 16)
				_, _ = that.ClientConn.Write([]byte(len64))
				_, _ = that.ClientConn.Write([]byte("\r\n"))
				_, _ = that.ClientConn.Write(body)
				_, _ = that.ClientConn.Write([]byte("\r\n0\r\n\r\n"))
			} else {
				_, _ = that.ClientConn.Write(body)
			}
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
					return
				}
			}
		}
	}
}
