package plugin

import (
	"browser-reptile/common/plugin/baidu"
	"browser-reptile/common/plugin/base"
	"browser-reptile/common/proxy"
)

var (
	Handle = make([]proxy.HttpHandle, 0)
)

func init() {
	Handle = append(Handle, &base.Base{})
	Handle = append(Handle, &baidu.Home{})
}
