package net_protocol

type Selector interface {
	Start()
	StartAndWait()
}

type INetProtocol interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	GetRemoteAddr() string
	GetLocalAddr() string
	Start(selector Selector)
}

type IListener interface {
	Listen(addr string, handler func(INetProtocol)) error
	Close() error
}
