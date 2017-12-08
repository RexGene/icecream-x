package net_protocol

import (
	"net"
)

type INetProtocol interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	GetRemoteAddr() net.Addr
}

type IListener interface {
	Listen(addr string, handler func(INetProtocol)) error
	Close() error
}
