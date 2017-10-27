package net_protocol

type INetProtocol interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type IListener interface {
	Listen(addr string, handler func(INetProtocol)) error
	Close() error
}
