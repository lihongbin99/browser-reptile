package main

import (
	"browser-reptile/config"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// 连接ID
	id     = 0
	idLock = sync.Mutex{}
)

func getId() int {
	idLock.Lock()
	defer idLock.Unlock()
	id++
	return id
}

func IsIpv4(hostname string) bool {
	sp := strings.Split(hostname, ".")
	if len(sp) != 4 {
		return false
	}
	for i := 0; i < len(sp); i++ {
		if _, err := strconv.Atoi(sp[i]); err != nil {
			return false
		}
	}
	return true
}

func WildcardDomain(hostname string) string {
	sp := strings.Split(hostname, ".")
	if len(sp) <= 2 {
		return hostname
	}

	domain := "*"
	for i := 1; i < len(sp); i++ {
		domain += "." + sp[i]
	}

	return domain
}

func Tunnel(sockId int, clientConn, serverConn net.Conn) {
	defer clientConn.Close()
	defer serverConn.Close()

	errCh := make(chan error)
	go Copy(sockId, true, clientConn, serverConn, errCh)
	go Copy(sockId, false, serverConn, clientConn, errCh)

	_ = <-errCh
	_ = <-errCh
}

func Copy(sockId int, isClient bool, src io.Reader, dst io.Writer, errCh chan error) {
	direction := ""
	if isClient {
		direction = ">>>"
	} else {
		direction = "<<<"
	}

	buf := make([]byte, 64*1024)
	for {
		readLen, err := src.Read(buf)
		if err != nil {
			WriteData(sockId, fmt.Sprintf("%s%s%s\n[%s]\n\n", direction, time.Now().String()[:19], direction, err.Error()))
			errCh <- err
			return
		}

		// 日志
		WriteData(sockId, fmt.Sprintf("%s%s%s\n[%s]\n\n", direction, time.Now().String()[:19], direction, string(buf[:readLen])))

		dst.Write(buf[:readLen])
	}
}

func WriteData(sockId int, data string) {
	if config.DataTo == config.DataToFile {
		fileName := fmt.Sprintf("%s%d.bin", config.DataDir, sockId)
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
		if err == nil {
			file.Write([]byte(data))
			file.Close()
		}
	} else if config.DataTo == config.DataToConsole {
		fmt.Print(data)
	}
}
