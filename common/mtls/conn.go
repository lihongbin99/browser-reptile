package mtls

import (
	"crypto/tls"
	"net"
)

type Conn struct {
	net.Conn
	isClient  bool
	Handshake func() error
	config    *tls.Config

	clientHello *ClientHelloMsg
	serverHello *ServerHelloMsg
}

func (c *Conn) GetClientHello() (*ClientHelloMsg, error) {
	if c.clientHello != nil {
		return c.clientHello, nil
	}
	// TODO 开始读取 Client Hello
	return nil, nil
}

func (c *Conn) SetServerHello(serverHello *ServerHelloMsg) {

}

func (c *Conn) serverHandshake() error {
	//clientHello, err := c.readClientHello()
	//if err != nil {
	//	return err
	//}
	//
	//if c.vers == VersionTLS13 {
	//	hs := serverHandshakeStateTLS13{
	//		c:           c,
	//		ctx:         ctx,
	//		clientHello: clientHello,
	//	}
	//	return hs.handshake()
	//}
	//
	//hs := serverHandshakeState{
	//	c:           c,
	//	ctx:         ctx,
	//	clientHello: clientHello,
	//}
	//return hs.handshake()
	return nil
}

func (c *Conn) GetServerHello() (*ServerHelloMsg, error) {
	return nil, nil
}

func (c *Conn) SetClientHello(clientHello *ClientHelloMsg) {
	//
}

func (c *Conn) clientHandshake() (err error) {
	//if c.config == nil {
	//	c.config = defaultConfig()
	//}
	//
	//// This may be a renegotiation handshake, in which case some fields
	//// need to be reset.
	//c.didResume = false
	//
	//hello, ecdheParams, err := c.makeClientHello()
	//if err != nil {
	//	return err
	//}
	//c.serverName = hello.serverName
	//
	//cacheKey, session, earlySecret, binderKey := c.loadSession(hello)
	//if cacheKey != "" && session != nil {
	//	defer func() {
	//		// If we got a handshake failure when resuming a session, throw away
	//		// the session ticket. See RFC 5077, Section 3.2.
	//		//
	//		// RFC 8446 makes no mention of dropping tickets on failure, but it
	//		// does require servers to abort on invalid binders, so we need to
	//		// delete tickets to recover from a corrupted PSK.
	//		if err != nil {
	//			c.config.ClientSessionCache.Put(cacheKey, nil)
	//		}
	//	}()
	//}
	//
	//if _, err := c.writeRecord(recordTypeHandshake, hello.marshal()); err != nil {
	//	return err
	//}
	//
	//msg, err := c.readHandshake()
	//if err != nil {
	//	return err
	//}
	//
	//serverHello, ok := msg.(*serverHelloMsg)
	//if !ok {
	//	c.sendAlert(alertUnexpectedMessage)
	//	return unexpectedMessageError(serverHello, msg)
	//}
	//
	//if err := c.pickTLSVersion(serverHello); err != nil {
	//	return err
	//}
	//
	//// If we are negotiating a protocol version that's lower than what we
	//// support, check for the server downgrade canaries.
	//// See RFC 8446, Section 4.1.3.
	//maxVers := c.config.maxSupportedVersion(roleClient)
	//tls12Downgrade := string(serverHello.random[24:]) == downgradeCanaryTLS12
	//tls11Downgrade := string(serverHello.random[24:]) == downgradeCanaryTLS11
	//if maxVers == VersionTLS13 && c.vers <= VersionTLS12 && (tls12Downgrade || tls11Downgrade) ||
	//	maxVers == VersionTLS12 && c.vers <= VersionTLS11 && tls11Downgrade {
	//	c.sendAlert(alertIllegalParameter)
	//	return errors.New("tls: downgrade attempt detected, possibly due to a MitM attack or a broken middlebox")
	//}
	//
	//if c.vers == VersionTLS13 {
	//	hs := &clientHandshakeStateTLS13{
	//		c:           c,
	//		ctx:         ctx,
	//		serverHello: serverHello,
	//		hello:       hello,
	//		ecdheParams: ecdheParams,
	//		session:     session,
	//		earlySecret: earlySecret,
	//		binderKey:   binderKey,
	//	}
	//
	//	// In TLS 1.3, session tickets are delivered after the handshake.
	//	return hs.handshake()
	//}
	//
	//hs := &clientHandshakeState{
	//	c:           c,
	//	ctx:         ctx,
	//	serverHello: serverHello,
	//	hello:       hello,
	//	session:     session,
	//}
	//
	//if err := hs.handshake(); err != nil {
	//	return err
	//}
	//
	//// If we had a successful handshake and hs.session is different from
	//// the one already cached - cache a new one.
	//if cacheKey != "" && hs.session != nil && session != hs.session {
	//	c.config.ClientSessionCache.Put(cacheKey, hs.session)
	//}

	return nil
}
