package parser

import (
	"encoding/json"
	"github.com/RexGene/icecreamx/proxy"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
)

type JSonParser struct {
	sync.RWMutex
	handlerMap map[uint]*handleNode
}

func NewJSonParser() *JSonParser {
	return &JSonParser{
		handlerMap: make(map[uint]*handleNode),
	}
}

func (self *JSonParser) Register(
	id uint,
	makeFunc func() proto.Message,
	handleFunc func(*proxy.NetProxy, proto.Message)) error {

	if handleFunc == nil || makeFunc == nil {
		return ErrHandleFuncIsNil
	}

	self.Lock()
	defer self.Unlock()

	_, isOK := self.handlerMap[id]
	if isOK {
		return ErrCmdIdAlreadyExist
	}

	self.handlerMap[id] = &handleNode{
		makeFunc:   makeFunc,
		handleFunc: handleFunc,
	}
	return nil
}

func (self *JSonParser) Unregister(id uint) {
	self.Lock()
	self.handlerMap[id] = nil
	self.Unlock()
}

func (self *JSonParser) ParseAndHandle(client *proxy.NetProxy, id uint, data []byte) error {
	self.RLock()
	defer self.RUnlock()

	node := self.handlerMap[id]
	if node == nil {
		return ErrCmdIdNotFound
	}

	log.Println("[!]", string(data))
	msg := node.makeFunc()
	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	node.handleFunc(client, msg)
	return nil
}

func (self *JSonParser) Marshal(msg proto.Message) ([]byte, error) {
	return json.Marshal(msg)
}
