package main

import (
	"browser-reptile/common/config"
	"browser-reptile/common/plugin"
	_ "browser-reptile/common/plugin/baidu"
	"browser-reptile/common/proxy"
	"browser-reptile/common/socks"
	"browser-reptile/web"
	"fmt"
	"net"
	"sync"
)

var (
	id     = 0
	idLock = sync.Mutex{}
)

func main() {

	listenAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", config.ListenPort))
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("start success")
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			if err != nil {
				panic(err)
			}
		}

		go doMain(conn)
	}
}

func doMain(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	sockId := getId()

	buf := make([]byte, 64*1024)

	socksProxy, local, err := socks.Parse(conn, buf)
	if err != nil {
		fmt.Printf("socks parse error: %v\n", err)
		return
	}

	if local {
		web.Main(conn)
		return
	}

	defer func() {
		_ = socksProxy.ServerConn.Close()
	}()

	//fmt.Printf("socks proxy %s:%d success\n", socksProxy.Host, socksProxy.Port)

	waitChan := make(chan uint8)
	if socksProxy.NeedParseHttp {
		httpProxy := proxy.NewHttpProxy(socksProxy, plugin.Handle, waitChan)
		go httpProxy.HttpRequest()
		go httpProxy.HttpResponse()
	} else {
		go proxy.Transfer(sockId, socksProxy.ClientConn, socksProxy.ServerConn, waitChan, true)
		go proxy.Transfer(sockId, socksProxy.ServerConn, socksProxy.ClientConn, waitChan, false)
	}

	_ = <-waitChan
	_ = <-waitChan
}

func getId() int {
	idLock.Lock()
	defer idLock.Unlock()
	id++
	return id
}
