package main

import (
	"browser-reptile/config"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/lihongbin99/btls"
	"github.com/lihongbin99/log"
	"github.com/lihongbin99/utils"

	tls "github.com/refraction-networking/utls"
)

var (
	CertMap = make(map[string]uint8)
)

func Tls(clientConn, serverConn net.Conn, hostname string) (net.Conn, net.Conn, error) {
	clientPeepConn := &utils.PeepIo{Conn: clientConn}

	// 偷窥是否 tls 协议连接
	buf, err := clientPeepConn.PeepN(1)
	if err != nil {
		log.Debug("偷窥是否 tls 协议连接失败:", err)
		return nil, nil, err
	}

	if buf[0] != 22 {
		return clientPeepConn, serverConn, nil
	}

	// 偷窥 ClientHello
	clientHello, err := PeepClientHello(clientPeepConn)
	if err != nil {
		log.Debug("偷窥 tls 协议失败:", err)
		return nil, nil, err
	}

	// 与远程 tls 连接
	utlsConn, err := TlsDialServer(clientHello, serverConn)
	if err != nil {
		log.Debug("与服务器进行 tls 握手失败:", err)
		return nil, nil, err
	}

	// 解析远程 tls 参数
	Protocol := utlsConn.ConnectionState().NegotiatedProtocol

	// 校验证书
	certDomain := hostname
	if !IsIpv4(certDomain) {
		certDomain = WildcardDomain(certDomain)
	}
	if _, exist := CertMap[certDomain]; !exist {
		cmd := exec.Command(config.CertificateUtils, certDomain)
		if err := cmd.Run(); err != nil {
			return nil, nil, fmt.Errorf("创建证书失败: %s", certDomain)
		}
		CertMap[certDomain] = 1
	}
	certDomain = strings.ReplaceAll(certDomain, "*", "x")

	// 构建 ServerHello
	cert, err := tls.LoadX509KeyPair(
		config.CertificatePath+"CertificateUtils_"+certDomain+".crt",
		config.CertificatePath+"CertificateUtils_"+certDomain+".key")
	if err != nil {
		log.Debug("加载证书失败:", err)
		return nil, nil, err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}, NextProtos: []string{Protocol}}

	// 与本地 tls 连接
	localTlsConn := tls.Server(clientPeepConn, tlsConfig)
	if err := localTlsConn.Handshake(); err != nil {
		log.Debug("与浏览器进行 tls 握手失败:", err)
		return nil, nil, err
	}

	log.Debug("解析", hostname, " tls 成功")
	return localTlsConn, utlsConn, nil
}

func PeepClientHello(clientConn *utils.PeepIo) (*btls.HandshakeClientHello, error) {
	btlsConn := btls.TLSConn{}
	protocol, err := clientConn.PeepN(8)
	if err != nil {
		return nil, err
	}

	clientHelloBuf, err := clientConn.PeepN(int(protocol[5])<<16 | int(protocol[6])<<8 | int(protocol[7]))
	if err != nil {
		return nil, err
	}

	return btlsConn.ParseHandshakeClientHello(clientHelloBuf)
}

func TlsDialServer(clientHello *btls.HandshakeClientHello, serverConn net.Conn) (*tls.UConn, error) {
	// 获取 hostname
	hostname := ""
	for i := 0; i < len(clientHello.Extensions); i++ {
		switch extension := clientHello.Extensions[i].(type) {
		case *btls.ExtensionServerName:
			if len(extension.ServerNameList) > 0 {
				hostname = extension.ServerNameList[0].ServerName
			}
			break
		}
	}

	// 构建 ClientHello
	config := tls.Config{ServerName: hostname}
	uTlsConn := tls.UClient(serverConn, &config, tls.HelloCustom)
	utlsClientHello, err := MakeClientHello(clientHello)
	if err != nil {
		return nil, err
	}

	// 与服务器进行tls连接
	if err := uTlsConn.ApplyPreset(utlsClientHello); err != nil {
		return nil, err
	}
	if err := uTlsConn.Handshake(); err != nil {
		return nil, err
	}

	return uTlsConn, nil
}

func MakeClientHello(clientHello *btls.HandshakeClientHello) (*tls.ClientHelloSpec, error) {
	var buf []byte
	var err error
	utlsClientHello := tls.ClientHelloSpec{
		TLSVersMax:         tls.VersionTLS13,
		TLSVersMin:         tls.VersionTLS10,
		CompressionMethods: []uint8{0},
		GetSessionID:       nil,
	}

	CipherSuites := make([]uint16, 0)
	for i := 0; i < len(clientHello.CipherSuite); i++ {
		if btls.IsGREASE(uint16(clientHello.CipherSuite[i])) {
			CipherSuites = append(CipherSuites, tls.GREASE_PLACEHOLDER)
		} else {
			CipherSuites = append(CipherSuites, uint16(clientHello.CipherSuite[i]))
		}
	}
	Extensions := make([]tls.TLSExtension, 0)
	for i := 0; i < len(clientHello.Extensions); i++ {
		switch extension := clientHello.Extensions[i].(type) {
		case *btls.ExtensionGREASE:
			Extensions = append(Extensions, &tls.UtlsGREASEExtension{Value: tls.GREASE_PLACEHOLDER, Body: extension.Data})
		case *btls.ExtensionServerName:
			hostname := ""
			if len(extension.ServerNameList) > 0 {
				hostname = extension.ServerNameList[0].ServerName
			}
			Extensions = append(Extensions, &tls.SNIExtension{ServerName: hostname})
		case *btls.ExtensionStatusRequest:
			Extensions = append(Extensions, &tls.StatusRequestExtension{})
		case *btls.ExtensionSupportedGroups:
			Curves := make([]tls.CurveID, 0)
			for i := 0; i < len(extension.SupportedGroups); i++ {
				if btls.IsGREASE(uint16(extension.SupportedGroups[i])) {
					Curves = append(Curves, tls.GREASE_PLACEHOLDER)
				} else {
					Curves = append(Curves, tls.CurveID(extension.SupportedGroups[i]))
				}
			}
			Extensions = append(Extensions, &tls.SupportedCurvesExtension{Curves: Curves})
		case *btls.ExtensionEcPointFormats:
			Extensions = append(Extensions, &tls.SupportedPointsExtension{SupportedPoints: []byte{0}})
		case *btls.ExtensionSignatureAlgorithms:
			SupportedSignatureAlgorithms := make([]tls.SignatureScheme, 0)
			for i := 0; i < len(extension.SignatureHashAlgorithms); i++ {
				SupportedSignatureAlgorithms = append(SupportedSignatureAlgorithms, tls.SignatureScheme(extension.SignatureHashAlgorithms[i]))
			}
			Extensions = append(Extensions, &tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: SupportedSignatureAlgorithms})
		case *btls.ExtensionAppclicationLayerProtocol:
			alpnProtocols := make([]string, 0)
			for i := 0; i < len(extension.ALPNProtocols); i++ {
				alpnProtocols = append(alpnProtocols, extension.ALPNProtocols[i].ALPNNextProtocol)
			}
			Extensions = append(Extensions, &tls.ALPNExtension{AlpnProtocols: alpnProtocols})
		case *btls.ExtensionSignedCertificateTimestamp:
			Extensions = append(Extensions, &tls.SCTExtension{})
		case *btls.ExtensionPadding:
			Extensions = append(Extensions, &tls.UtlsPaddingExtension{PaddingLen: len(extension.PaddingData), WillPad: true})
		case *btls.ExtensionExtendedMasterSecret:
			Extensions = append(Extensions, &tls.UtlsExtendedMasterSecretExtension{})
		case *btls.ExtensionCompressCertificate:
			Algorithms := make([]tls.CertCompressionAlgo, 0)
			for i := 0; i < len(extension.Algorithms); i++ {
				Algorithms = append(Algorithms, tls.CertCompressionAlgo(extension.Algorithms[i]))
			}
			Extensions = append(Extensions, &tls.UtlsCompressCertExtension{Algorithms: Algorithms})
		case *btls.ExtensionSessionTicket:
			Extensions = append(Extensions, &tls.SessionTicketExtension{})
		case *btls.ExtensionSupportedVersion:
			Versions := make([]uint16, 0)
			for i := 0; i < len(extension.SupportedVersions); i++ {
				Versions = append(Versions, uint16(extension.SupportedVersions[i]))
			}
			Extensions = append(Extensions, &tls.SupportedVersionsExtension{Versions: Versions})
		case *btls.ExtensionPskKeyExchangeModes:
			Modes := make([]uint8, 0)
			for i := 0; i < len(extension.PSKKeyExchangeMode); i++ {
				Modes = append(Modes, uint8(extension.PSKKeyExchangeMode[i]))
			}
			Extensions = append(Extensions, &tls.PSKKeyExchangeModesExtension{Modes: Modes})
		case *btls.ExtensionKeyShare:
			KeyShares := make([]tls.KeyShare, 0)
			for i := 0; i < len(extension.KeyShareEntrys); i++ {
				if btls.IsGREASE(uint16(extension.KeyShareEntrys[i].Group)) {
					KeyShares = append(KeyShares, tls.KeyShare{Group: tls.GREASE_PLACEHOLDER, Data: []byte{0}})
				} else {
					KeyShares = append(KeyShares, tls.KeyShare{Group: tls.CurveID(extension.KeyShareEntrys[i].Group)})
				}
			}
			Extensions = append(Extensions, &tls.KeyShareExtension{KeyShares: KeyShares})
		case *btls.ExtensionApplicationSettings:
			SupportedProtocols := make([]string, 0)
			for i := 0; i < len(extension.SupportedALPNList); i++ {
				SupportedProtocols = append(SupportedProtocols, extension.SupportedALPNList[i].SupportedALPN)
			}
			Extensions = append(Extensions, &tls.ApplicationSettingsExtension{SupportedProtocols: SupportedProtocols})
		case *btls.ExtendedRenegotiationInfo:
			Extensions = append(Extensions, &tls.RenegotiationInfoExtension{Renegotiation: 0})
		case *btls.ExtensionPreSharedKey:
			// utls 没有实现 PreSharedKey
			buf, err = extension.ToBuf(nil)
			if err != nil {
				return nil, err
			}
			Extensions = append(Extensions, &tls.GenericExtension{Id: uint16(btls.EXTENSION_PRE_SHARED_KEY), Data: buf[4:]})
		default:
			return nil, fmt.Errorf("没有解析Extensions")
		}
	}

	utlsClientHello.CipherSuites = CipherSuites
	utlsClientHello.Extensions = Extensions

	return &utlsClientHello, nil
}
