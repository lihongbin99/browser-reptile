package base

import "browser-reptile/common/proxy"

type Default struct{}

func (that *Default) NeedParse(host string, port uint16) bool {
	return false
}

func (that *Default) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Default) RequestParseHeader(proxy *proxy.HttpProxy) {
	//for headerName, headerValues := range proxy.RequestHeader { }
}

func (that *Default) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Default) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *Default) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Default) ResponseParseHeader(proxy *proxy.HttpProxy) {
	//for headerName, headerValues := range proxy.ResponseHeader { }
}

func (that *Default) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Default) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}
