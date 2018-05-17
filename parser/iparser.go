package parser

import (
	"github.com/RexGene/icecreamx/proxy"
	"github.com/golang/protobuf/proto"
)

type IParser interface {
	Register(id uint, makeFunc func() proto.Message, handleFunc func(*proxy.NetProxy, proto.Message)) error
	Unregister(id uint)
	ParseAndHandle(client *proxy.NetProxy, id uint, data []byte) error
	Marshal(msg proto.Message) ([]byte, error)
}
