package proxy

import (
	"net"
)

type ServerProxy struct {
	NetProxy
}

func NewServerProxy(conn net.Conn, parser IParser,
	netProtocol INetProtocol) *ServerProxy {
	proxy := &ServerProxy{}
	proxy.isRunning = false
	proxy.conn = conn
	proxy.parser = parser
	proxy.netProtocol = netProtocol

	return proxy
}
