package proxy

import (
	"fmt"
	"io"
	"net"
)

type Proxy struct {
	ClientConn net.Conn
	ServerConn net.Conn

	Host string
	Port uint16

	NeedParse bool
}

func Transfer(r net.Conn, w net.Conn, waitChan chan uint8) {
	transfer(r, w, waitChan, false)
}

func TransferShowContent(r net.Conn, w net.Conn, waitChan chan uint8) {
	transfer(r, w, waitChan, true)
}

func transfer(r net.Conn, w net.Conn, waitChan chan uint8, flag bool) {
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
		if flag {
			fmt.Printf("proxy: %s\n", string(buf[:readLen]))
		}

		if _, err = w.Write(buf[:readLen]); err != nil {
			return
		}

		if err == io.EOF {
			break
		}
	}
}
