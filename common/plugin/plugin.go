package plugin

import (
	"browser-reptile/common/plugin/baidu"
	"browser-reptile/common/plugin/lhb13520"
	"browser-reptile/common/plugin/test789590"
	"browser-reptile/common/proxy"
)

var (
	Handle = make([]proxy.HttpHandle, 0)
)

func init() {
	Handle = append(Handle, &baidu.Home{})
	Handle = append(Handle, &lhb13520.All{})
	Handle = append(Handle, &test789590.All{})
}
