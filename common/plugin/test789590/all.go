package test789590

import (
	"browser-reptile/common/proxy"
	"fmt"
	"strings"
)

type All struct{}

func (that *All) NeedParse(host string, port uint16) bool {
	return strings.HasSuffix(host, "test789590.com")
}

func (that *All) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *All) RequestParseHeader(proxy *proxy.HttpProxy) {
}

func (that *All) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *All) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *All) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *All) ResponseParseHeader(proxy *proxy.HttpProxy) {
}

func (that *All) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *All) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	fmt.Println(string(body))
	return body
}
