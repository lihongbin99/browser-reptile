package base

import "browser-reptile/common/proxy"

type Default struct{}

func (that *Default) NeedParse(host string, port uint16) bool {
	return false
}

func (that *Default) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Default) RequestParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Default) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Default) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *Default) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Default) ResponseParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Default) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false
}

func (that *Default) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}
