package proxy

import ()

type ClientProxy struct {
	NetProxy
}

func NewClientProxy(netProtocol INetProtocol, recvicer ICloseNotifyRecvicer,
	parser IParser) *ClientProxy {
	proxy := &ClientProxy{}

	proxy.isRunning = false
	proxy.parser = parser
	proxy.recvicer = recvicer
	proxy.netProtocol = netProtocol

	return proxy
}
