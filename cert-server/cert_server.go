package main

import (
	"browser-reptile/common/config"
	"browser-reptile/web"
	"fmt"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", config.ListenPort))
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			break
		}
		go web.Main(conn)
	}
}
