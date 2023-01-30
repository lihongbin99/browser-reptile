package main

import (
	"browser-reptile/config"
	"fmt"
	"net"

	"github.com/lihongbin99/log"
	"github.com/lihongbin99/proxy_server"
	"github.com/lihongbin99/utils"
)

func main() {
	var err error
	config.Listener, err = proxy_server.Listen(fmt.Sprintf(":%d", config.ListenPort))
	if err != nil {
		panic(err)
	}

	log.Info("server start success")
	for {
		proxy, err := config.Listener.Accept()
		if err != nil {
			break
		}

		go doMain(proxy)
	}

	config.Listener.Close()
	log.Info("server exit success")
}

func doMain(proxy *proxy_server.Proxy) {
	defer proxy.Close()
	var err error

	// 获取连接ID
	sockId := getId()

	var clientConn, serverConn net.Conn

	// 解析 tls
	if config.ParseTLS {
		if clientConn, serverConn, err = Tls(proxy.ClientConn, proxy.ServerConn, proxy.HOSTNAME); err != nil {
			log.Warn("sockId:", sockId, " ---> ", proxy.HOSTNAME, "tls连接错误:", err)
			return
		}
	} else {
		clientConn = proxy.ClientConn
		serverConn = proxy.ServerConn
	}

	// 解析 HTTP
	if config.ParseHTTP {
		peepConn := utils.PeepIo{Conn: clientConn}
		protocol, err := peepConn.PeepN(1)
		if err != nil {
			return
		}
		if protocol[0] == 22 { // 无法解析 tls
			Tunnel(sockId, &peepConn, serverConn)
			return
		}
		Http(sockId, &peepConn, serverConn)
	} else {
		Tunnel(sockId, clientConn, serverConn)
	}
}
