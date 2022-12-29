package mtls

import (
	"browser-reptile/common/config"
	"crypto/tls"
	"net"
)

func InitTLS(clientConn net.Conn, serverAddr string) (*Conn, *Conn, error) {
	// 连接服务器
	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, nil, err
	}

	// 初始化TLS连接
	clientTLS := server(clientConn, config.TlsConfig)
	serverTLS := client(serverConn)

	// 1. 获取 Client Hello
	clientHello, err := clientTLS.GetClientHello()
	if err != nil {
		return nil, nil, err
	}

	// 2. 仿造 Client Hello 中的全部参数
	serverTLS.SetClientHello(clientHello)

	// 3. 与服务器握手
	if err = serverTLS.Handshake(); err != nil {
		return nil, nil, err
	}

	// 4. 获取 Server Hello
	serverHello, err := serverTLS.GetServerHello()
	if err != nil {
		return nil, nil, err
	}

	// 5. 仿造 Server Hello 中的指定几个参数
	clientTLS.SetServerHello(serverHello)

	// 6. 与浏览器握手
	if err = clientTLS.Handshake(); err != nil {
		return nil, nil, err
	}

	return clientTLS, serverTLS, nil
}

func server(conn net.Conn, config *tls.Config) *Conn {
	c := &Conn{
		Conn:   conn,
		config: config,
	}
	c.Handshake = c.serverHandshake
	return c
}

func client(conn net.Conn) *Conn {
	c := &Conn{
		Conn:     conn,
		isClient: true,
	}
	c.Handshake = c.clientHandshake
	return c
}
