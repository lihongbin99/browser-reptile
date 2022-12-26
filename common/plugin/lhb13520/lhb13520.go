package lhb13520

import (
	"browser-reptile/common/proxy"
	"strings"
)

type Home struct{}

func (that *Home) NeedParse(host string, port uint16) bool {
	return strings.HasSuffix(host, "lhb13520.com")
}

func (that *Home) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) RequestParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Home) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *Home) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) ResponseParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Home) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}
