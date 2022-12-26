package baidu

import (
	"browser-reptile/common/proxy"
	"strings"
)

type Home struct{}

func (that *Home) NeedParse(host string, port uint16) bool {
	return strings.HasSuffix(host, "baidu.com")
}

func (that *Home) RequestNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) RequestParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Home) RequestNeedParseBody(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) RequestParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	return body
}

func (that *Home) ResponseNeedParseHeader(proxy *proxy.HttpProxy) bool {
	return false && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) ResponseParseHeader(proxy *proxy.HttpProxy) {
}

func (that *Home) ResponseNeedParseBody(proxy *proxy.HttpProxy) bool {
	return true && that.NeedParse(proxy.Host, proxy.Port)
}

func (that *Home) ResponseParseBody(proxy *proxy.HttpProxy, body []byte) []byte {
	body = []byte(strings.ReplaceAll(string(body), "id=\"index-kw\"", "id=\"index-kw\"  placeholder=\"success\""))
	return []byte(strings.ReplaceAll(string(body), "id=\"kw\"", "id=\"kw\"  placeholder=\"success\""))
}
