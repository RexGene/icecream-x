package proxy

import (
	"github.com/RexGene/icecreamx/net_protocol"
)

type ClientProxy struct {
	NetProxy
}

func NewClientProxy(netProtocol net_protocol.INetProtocol, recvicer ICloseNotifyRecvicer,
	parser IParser) *ClientProxy {
	proxy := &ClientProxy{}

	proxy.isRunning = false
	proxy.parser = parser
	proxy.recvicer = recvicer
	proxy.netProtocol = netProtocol

	return proxy
}
