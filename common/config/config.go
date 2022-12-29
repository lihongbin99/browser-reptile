package config

import (
	"crypto/tls"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"time"
)

var (
	configPath = "config.ini"
	iniConfig  *ini.File

	ListenPort = 0

	TlsCertFile = ""
	TlsKeyFile  = ""
	TlsConfig   *tls.Config

	CommonConfig CommonConfigStruct

	LogDir string
)

type CommonConfigStruct struct {
	ParseTLS  bool
	ParseHTTP bool
	LogTo     LogToType
}

type LogToType int8

const (
	_ LogToType = iota
	LogToConsole
	LogToFile
	LogToNull
)

func init() {
	flag.StringVar(&configPath, "c", configPath, "config file")
	flag.Parse()

	var err error
	iniConfig, err = ini.Load(configPath)
	if err != nil {
		panic(err)
	}

	ListenPort, err = iniConfig.Section("server").Key("listen_port").Int()
	if err != nil {
		fmt.Println("get listen port error")
		panic(err)
	}

	TlsCertFile = iniConfig.Section("tls").Key("crt_path").String()
	TlsKeyFile = iniConfig.Section("tls").Key("key_path").String()

	fmt.Println("config listen port:", ListenPort)
	fmt.Println("config tls cert:", TlsCertFile)
	fmt.Println("config tls key:", TlsKeyFile)

	cert, err := tls.LoadX509KeyPair(TlsCertFile, TlsKeyFile)
	if err != nil {
		panic(err)
	}
	TlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

	CommonConfig.ParseTLS = iniConfig.Section("common").Key("parse_tls").String() == "true"
	CommonConfig.ParseHTTP = iniConfig.Section("common").Key("parse_http").String() == "true"
	logTo := iniConfig.Section("common").Key("log_to").String()
	if logTo == "console" {
		CommonConfig.LogTo = LogToConsole
	} else if logTo == "file" {
		LogDir = "log" + string(os.PathSeparator) + strings.ReplaceAll(time.Now().String()[:19], ":", "-")
		_ = os.MkdirAll(LogDir, 0700)
		CommonConfig.LogTo = LogToFile
	} else if logTo == "null" {
		CommonConfig.LogTo = LogToNull
	} else {
		panic("config log_to error: " + logTo)
	}

	fmt.Println("config parse tls :", CommonConfig.ParseTLS)
	fmt.Println("config parse http:", CommonConfig.ParseHTTP)
	fmt.Println("config log to:", logTo)
}
