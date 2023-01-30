package main

import (
	"fmt"
	"net"

	"github.com/lihongbin99/log"
)

type Req struct {
	Method string
	Path   string
	Proto  string
	Header map[string][]string
	Body   []byte
}

type Rep struct {
	Proto  string
	Status string
	Header map[string][]string
	Body   []byte
}

func Http(sockId int, clientConn, serverConn net.Conn) {
	for {
		req, err := ReadRequest(clientConn)
		if err != nil {
			break
		}

		// TODO 修改 Request
		reqHeader := "Request[\n"
		for k, v := range req.Header {
			reqHeader += fmt.Sprintf("%s: %v\n", k, v)
		}
		reqHeader += "]"
		log.Info(reqHeader)

		if err := SendRequest(req, serverConn); err != nil {
			break
		}

		rep, err := ReadResponse(serverConn)
		if err != nil {
			break
		}

		// TOOD 修改 Response
		repHeader := "Response[\n"
		for k, v := range rep.Header {
			repHeader += fmt.Sprintf("%s: %v\n", k, v)
		}
		repHeader += "]"
		log.Info(repHeader)

		if err := SendResponse(rep, clientConn); err != nil {
			break
		}
	}
}
