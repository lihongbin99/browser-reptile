package proxy

import (
	"browser-reptile/common/config"
	"browser-reptile/common/utils"
	"io"
	"net"
)

type Proxy struct {
	ClientConn net.Conn
	ServerConn net.Conn

	Host string
	Port uint16

	NeedParseHttp bool
}

func Transfer(sockId int, r net.Conn, w net.Conn, waitChan chan uint8, isTo bool) {
	defer func() {
		_ = r.Close()
		_ = w.Close()
		waitChan <- 0
	}()
	buf := make([]byte, 64*1024)
	for {
		readLen, err := r.Read(buf)
		if (err != nil && err != io.EOF) || readLen <= 0 {
			return
		}

		// 保存数据
		if config.CommonConfig.LogTo != config.LogToNull {
			utils.SaveData(sockId, buf[:readLen], isTo)
		}

		if _, err = w.Write(buf[:readLen]); err != nil {
			return
		}

		if err == io.EOF {
			break
		}
	}
}
