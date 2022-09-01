package base

import (
	"browser-reptile/common/proxy"
)

type Base struct{}

func (that *Base) NeedParse(host string, port uint16) bool {
	return false
}

func (that *Base) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return true
}

func (that *Base) RequestParseHeader(proxy *proxy.HttpProxy) {
	// TODO 处理因为长连接引发的请求不断开问题
	values := make([]string, 1)
	values[0] = "close"
	proxy.RequestHeader["Connection"] = values
}

func (that *Base) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Base) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *Base) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Base) ResponseParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Base) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Base) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}
