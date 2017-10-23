package proxy

import (
	"net"
)

type ClientProxy struct {
	NetProxy
}

func NewClientProxy(conn net.Conn, recvicer ICloseNotifyRecvicer,
	parser IParser, netProtocol INetProtocol) *ClientProxy {
	proxy := &ClientProxy{}

	proxy.isRunning = false
	proxy.conn = conn
	proxy.parser = parser
	proxy.recvicer = recvicer
	proxy.netProtocol = netProtocol

	return proxy
}
