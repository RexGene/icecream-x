package icecreamx

import (
    "github.com/RexGene/icecream-x/proxy"
	"github.com/golang/protobuf/proto"
    "sync"
    "errors"
)

var (
    ErrCmdIdAlreadyExist = errors.New("<Handler> cmdId already exist")
    ErrCmdIdNotFound = errors.New("<Handle> cmdId not found")
    ErrHandleFuncIsNil = errors.New("<Handler> handler func is nil")
)

type handleNode struct {
    makeFunc func() proto.Message
    handleFunc func (*proxy.NetProxy, proto.Message)
}

type Handler struct {
    sync.RWMutex
	handlerMap   map[uint] *handleNode
}

func NewHandler() *Handler {
    return &Handler{
        handlerMap : make(map[uint] *handleNode),
    }
}

func (self *Handler) Register(
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

func (self *Handler) Unregister(id uint) {
    self.Lock()
    self.handlerMap[id] = nil
    self.Unlock()
}

func (self *Handler) ParseAndHandle(client *proxy.NetProxy, id uint, data []byte) error {
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

