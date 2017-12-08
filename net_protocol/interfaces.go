package net_protocol

import ()

type INetProtocol interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	GetRemoteAddr() string
	GetLocalAddr() string
}

type IListener interface {
	Listen(addr string, handler func(INetProtocol)) error
	Close() error
}
