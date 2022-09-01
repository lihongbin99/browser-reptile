package config

import (
	"crypto/tls"
	"flag"
	"fmt"
)

var (
	ListenPort = 8080

	TlsCertFile = "./tls/tls.crt"
	TlsKeyFile  = "./tls/tls.key"
	TlsConfig   *tls.Config
)

func init() {
	flag.IntVar(&ListenPort, "p", ListenPort, "listen port: 8080")
	flag.StringVar(&TlsCertFile, "c", TlsCertFile, "tls cert file")
	flag.StringVar(&TlsKeyFile, "k", TlsKeyFile, "tls key file")
	flag.Parse()

	fmt.Println("listen port:", ListenPort)
	fmt.Println("tls cert:", TlsCertFile)
	fmt.Println("tls key:", TlsKeyFile)

	cert, err := tls.LoadX509KeyPair(TlsCertFile, TlsKeyFile)
	if err != nil {
		panic(err)
	}
	TlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
}
