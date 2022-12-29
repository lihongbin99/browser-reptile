package socks

import (
	"browser-reptile/common/config"
	"browser-reptile/common/plugin"
	"browser-reptile/common/proxy"
	"browser-reptile/common/utils"
	"crypto/tls"
	"fmt"
	"net"
)

var (
	localhostAddressMap = make(map[string]uint8)
)

func init() {
	addressList, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addressList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localhostAddressMap[ipNet.IP.String()] = 127
			}
		}
	}
	localhostAddressMap["127.0.0.1"] = 127
	localhostAddressMap["localhost"] = 127
}

func Parse(conn net.Conn, buf []byte) (*proxy.Proxy, bool, error) {
	if err := step1(conn, buf); err != nil {
		return nil, false, err
	}
	return step2(conn, buf)
}

func step1(conn net.Conn, buf []byte) error {
	readLen, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("socks5 step1 read error: %v", err)
	}
	if readLen <= 0 {
		return fmt.Errorf("socks5 step1 length error: %v", readLen)
	}
	if readLen != 3 || buf[0] != 5 || buf[1] != 1 || buf[2] != 0 {
		return fmt.Errorf("socks5 step1 content error: %v", buf[:readLen])
	}
	buf[0] = 5
	buf[1] = 0
	_, err = conn.Write(buf[:2])
	return err
}

func step2(conn net.Conn, buf []byte) (*proxy.Proxy, bool, error) {
	readLen, err := conn.Read(buf)
	if err != nil {
		return nil, false, fmt.Errorf("socks5 step2 read error: %v", err)
	}
	if readLen < 4 {
		buf[1] = 1
		_, _ = conn.Write(buf[:readLen])
		return nil, false, fmt.Errorf("socks5 step2 length error: %v", readLen)
	}
	if buf[0] != 5 || buf[1] != 1 || buf[2] != 0 {
		buf[1] = 1
		_, _ = conn.Write(buf[:readLen])
		return nil, false, fmt.Errorf("socks5 step2 content error: %v", buf[:readLen])
	}

	host := ""
	switch buf[3] {
	case 1:
		host = fmt.Sprintf("%d.%d.%d.%d", buf[4], buf[5], buf[6], buf[7])
	case 3:
		host = string(buf[5 : 5+buf[4]])
	case 4:
		buf[1] = 2
		_, _ = conn.Write(buf[:readLen])
		return nil, false, fmt.Errorf("socks5 step2 not support ipv6: %v", buf[:readLen]) // TODO 暂时不支持 ipv6
	default:
		buf[1] = 2
		_, _ = conn.Write(buf[:readLen])
		return nil, false, fmt.Errorf("socks5 step2 content error: %v", buf[:readLen])
	}

	port := utils.ToUint16(buf[readLen-2 : readLen])

	// 判断是否回环地址
	if _, localhostAddress := localhostAddressMap[host]; localhostAddress {
		if int(port) == config.ListenPort {
			// 本地服务
			buf[1] = 0
			_, _ = conn.Write(buf[:readLen])
			return nil, true, nil
		}
		buf[1] = 4
		_, _ = conn.Write(buf[:readLen])
		return nil, false, fmt.Errorf("localhost address: %s:%d", host, port)
	}

	// 判断是否需要解析 http
	needParseHttp := false
	if config.CommonConfig.ParseHTTP {
		// 如果是443端口则必须解析tls的情况下才能解析http
		if port != 443 || config.CommonConfig.ParseTLS {
			for _, handle := range plugin.Handle {
				if handle.NeedParse(host, port) {
					needParseHttp = true
					break
				}
			}
		}
	}

	// 连接服务器
	clientConn := conn
	var serverConn net.Conn = nil
	if port == 443 && config.CommonConfig.ParseTLS {
		//clientConn, serverConn, err = mtls.InitTLS(conn, fmt.Sprintf("%s:%d", host, port))
		clientConn = tls.Server(conn, config.TlsConfig)
		serverConn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			buf[1] = 3
			_, _ = conn.Write(buf[:readLen])
			return nil, false, fmt.Errorf("socks5 step2 connect tls server error: %v", err)
		}
	} else {
		var serverAddr *net.TCPAddr = nil
		serverAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			buf[1] = 3
			_, _ = conn.Write(buf[:readLen])
			return nil, false, fmt.Errorf("socks5 step2 resolve server adderss error: %v", err)
		}
		serverConn, err = net.DialTCP("tcp", nil, serverAddr)
		if err != nil {
			buf[1] = 3
			_, _ = conn.Write(buf[:readLen])
			return nil, false, fmt.Errorf("socks5 step2 connect server error: %v", err)
		}
	}

	// 返回 socks 成功
	buf[1] = 0
	_, _ = conn.Write(buf[:readLen])

	result := &proxy.Proxy{
		ClientConn:    clientConn,
		ServerConn:    serverConn,
		Host:          host,
		Port:          port,
		NeedParseHttp: needParseHttp,
	}
	return result, false, nil
}
