package proxy

import (
	"github.com/RexGene/icecreamx/net_protocol"
	"net"
)

type ServerProxy struct {
	NetProxy
}

func NewServerProxy(conn net.Conn, parser IParser,
	netProtocol net_protocol.INetProtocol) *ServerProxy {
	proxy := &ServerProxy{}
	proxy.isRunning = false
	proxy.conn = conn
	proxy.parser = parser
	proxy.netProtocol = netProtocol

	return proxy
}
