package icecreamx

import (
	"github.com/RexGene/icecreamx/net_protocol"
	"github.com/RexGene/icecreamx/parser"
	"github.com/RexGene/icecreamx/proxy"
	"github.com/golang/protobuf/proto"
	"net"
)

type Client struct {
	serverProxy *proxy.ServerProxy
	parser      *parser.PbParser
}

func NewClient() *Client {
	return &Client{
		parser: parser.NewPbParser(),
	}
}

func (self *Client) Connect(addr string) (*proxy.ServerProxy, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	serverProxy := proxy.NewServerProxy(conn, self.parser,
		net_protocol.NewTcp(conn))

	serverProxy.Setup(proxy.NewDataBufferMaker(BUFFER_SIZE))
	serverProxy.Start()
	self.serverProxy = serverProxy

	return serverProxy, nil
}

func (self *Client) ConnectWithDataBuffer(addr string,
	dataBufferMaker *proxy.DataBufferMaker) (*proxy.ServerProxy, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	serverProxy := proxy.NewServerProxy(conn, self.parser,
		net_protocol.NewTcp(conn))
	serverProxy.Setup(proxy.NewDataBufferMaker(BUFFER_SIZE))
	serverProxy.Start()
	self.serverProxy = serverProxy

	return serverProxy, nil
}

func (self *Client) RegisterCommand(
	id uint,
	makeFunc func() proto.Message,
	handleFunc func(*proxy.NetProxy, proto.Message)) {
	self.parser.Register(id, makeFunc, handleFunc)
}

func (self *Client) UnregisterCommand(id uint) {
	self.parser.Unregister(id)
}
