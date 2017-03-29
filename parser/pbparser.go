package parser

import (
    "github.com/RexGene/icecreamx/proxy"
    "github.com/golang/protobuf/proto"
    "sync"
    "errors"
)

var (
    ErrCmdIdAlreadyExist = errors.New("<PbParser> cmdId already exist")
    ErrCmdIdNotFound = errors.New("<PbParser> cmdId not found")
    ErrHandleFuncIsNil = errors.New("<PbParser> handler func is nil")
)

type handleNode struct {
    makeFunc func() proto.Message
    handleFunc func (*proxy.NetProxy, proto.Message)
}

type PbParser struct {
    sync.RWMutex
    handlerMap   map[uint] *handleNode
}

func NewPbParser() *PbParser {
    return &PbParser{
        handlerMap : make(map[uint] *handleNode),
    }
}

func (self *PbParser) Register(
        id uint,
        makeFunc func() proto.Message,
        handleFunc func (*proxy.NetProxy, proto.Message)) error {

    if handleFunc == nil || makeFunc == nil {
        return ErrHandleFuncIsNil
    }

    self.Lock()
    defer self.Unlock()

    _, isOK := self.handlerMap[id]
    if isOK {
        return ErrCmdIdAlreadyExist
    }

    self.handlerMap[id] = &handleNode {
        makeFunc : makeFunc,
        handleFunc : handleFunc,
    }
    return nil
}

func (self *PbParser) Unregister(id uint) {
    self.Lock()
    self.handlerMap[id] = nil
    self.Unlock()
}

func (self *PbParser) ParseAndHandle(client *proxy.NetProxy, id uint, data []byte) error {
    self.RLock()
    defer self.RUnlock()

    node := self.handlerMap[id]
    if node == nil {
        return ErrCmdIdNotFound
    }

    msg := node.makeFunc()
    if err := proto.Unmarshal(data, msg); err != nil {
        return err
    }

    node.handleFunc(client, msg)
    return nil
}

