package proxy

import (
	"net"
)

type ServerProxy struct {
	NetProxy
}

func NewServerProxy(conn net.Conn, parser IParser) *ServerProxy {
	proxy := &ServerProxy{}
	proxy.isRunning = false
	proxy.conn = conn
	proxy.parser = parser

	return proxy
}
