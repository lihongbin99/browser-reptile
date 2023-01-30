package config

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/lihongbin99/log"
	"github.com/lihongbin99/proxy_server"
	"gopkg.in/ini.v1"
)

var (
	// 配置文件的路径
	configFilePath = "config.ini"

	// 服务器监听端口号
	ListenPort = 0
	Listener   *proxy_server.Listener

	// 是否解密 tls 数据
	ParseTLS bool

	// 是否解密 http 数据
	ParseHTTP bool

	// 日志打印位置
	DataTo  DataToType
	DataDir string

	// 证书工具
	CertificateUtils string
	CertificatePath  string
)

type DataToType int8

const (
	_ DataToType = iota
	DataToConsole
	DataToFile
	DataToNull
)

func init() {
	pwd, _ := os.Getwd()

	// 设置配置文件路径
	flag.StringVar(&configFilePath, "c", configFilePath, "config file")
	flag.Parse()

	// 读取配置文件
	iniConfig, err := ini.Load(configFilePath)
	if err != nil {
		panic(err)
	}

	// 解析服务器配置
	ListenPort, err = iniConfig.Section("server").Key("listen_port").Int()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	log.ChangeLevel(strings.ToLower(iniConfig.Section("common").Key("log_level").String()))

	// 解析项目配置
	ParseTLS = strings.ToLower(iniConfig.Section("common").Key("parse_tls").String()) == "true"
	ParseHTTP = strings.ToLower(iniConfig.Section("common").Key("parse_http").String()) == "true"
	dataTo := strings.ToLower(iniConfig.Section("common").Key("data_to").String())
	if dataTo == "console" {
		DataTo = DataToConsole
	} else if dataTo == "file" {
		DataDir = pwd + string(os.PathSeparator) + "log" + string(os.PathSeparator) + strings.ReplaceAll(time.Now().String()[:19], ":", "-") + string(os.PathSeparator)
		_ = os.MkdirAll(DataDir, 0700)
		DataTo = DataToFile
	} else if dataTo == "null" {
		DataTo = DataToNull
	} else {
		panic("日志路径错误: " + dataTo + " not in [console, file, null]")
	}

	log.Info("server listen port: ", ListenPort)
	log.Info("parse tls         : ", ParseTLS)
	log.Info("parse http        : ", ParseHTTP)
	log.Info("data to           : ", dataTo)

	// 证书工具
	CertificateUtils = pwd + string(os.PathSeparator) + "Certificate-Utils.exe"
	if _, err = os.Stat(CertificateUtils); err != nil {
		panic("未找到: " + CertificateUtils)
	}
	CertificatePath = pwd + string(os.PathSeparator) + "cert" + string(os.PathSeparator)
}
