package icecreamx

import (
    "net"
    "log"
    "sync"
    "github.com/RexGene/icecreamx/proxy"
    "github.com/RexGene/icecreamx/parser"
    "github.com/golang/protobuf/proto"
)

const (
    DEFUALT_CHANNEL_SIZE = 32
)

type Server struct {
    clientSetMutex  sync.RWMutex
    runningMutex    sync.RWMutex
    waitGroup       sync.WaitGroup
    isRunning       bool
    addr            string
    listener        net.Listener
    parser         *parser.PbParser
    clientSet       map[interface{}] bool
    chRemoveClient  chan interface{}
}

func NewServer(addr string) (*Server, error) {
    server := &Server {
        isRunning : false,
        listener : nil,
        parser : parser.NewPbParser(),
    }

    err := server.init(addr)
    if err != nil {
        return nil, err
    }

    return server, nil
}

func (self *Server) Start() {
    if ! self.isRunning {
        self.isRunning = true
        self.clientSet = make(map[interface{}] bool)
        self.chRemoveClient = make(chan interface{}, DEFUALT_CHANNEL_SIZE)

        go self.listen_execute()
        go self.eventloop_execute()
    }
}

func (self *Server) Stop() {
    runningMutex := self.runningMutex
    runningMutex.Lock()
    self.isRunning = false
    runningMutex.Unlock()

    if self.listener != nil {
        self.listener.Close()
    }

    close(self.chRemoveClient)

    clientSetMutex := self.clientSetMutex
    clientSetMutex.RLock()
    for clientProxy, _ := range self.clientSet {
        clientProxy.(*proxy.ClientProxy).Stop()
    }
    clientSetMutex.RUnlock()

    self.waitGroup.Wait()
}

func (self *Server) PushCloseNotify(v interface{}) {
    runningMutex := self.runningMutex
    runningMutex.RLock()
    if self.isRunning {
        self.chRemoveClient <- v
    }
    runningMutex.RUnlock()
}

func (self *Server) listen_execute() {
    self.waitGroup.Add(1)
    defer self.waitGroup.Done()

    listener, err := net.Listen("tcp", self.addr)
    if err != nil {
        log.Println("[-]", err)
        return
    }
    defer listener.Close()
    self.listener = listener

    log.Println("[*] listening...")
    for self.isRunning {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("[!]", err)
        }

        clientProxy := proxy.NewClientProxy(conn, self, self.parser)
        clientProxy.Start()

        clientSetMutex := self.clientSetMutex
        clientSetMutex.Lock()
        self.clientSet[clientProxy] = true
        clientSetMutex.Unlock()
    }
}

func (self *Server) eventloop_execute() {
    self.waitGroup.Add(1)
    defer self.waitGroup.Done()

    for self.isRunning {
        select {
            case clientProxy, isOK := <-self.chRemoveClient:
                if isOK {
                    clientSetMutex := self.clientSetMutex
                    clientSetMutex.Lock()
                    delete(self.clientSet, clientProxy)
                    clientSetMutex.Unlock()
                }
        }
    }
}

func (self *Server) init(addr string) error {
    self.addr = addr
    return nil
}

func (self *Server) getConnectionCount() uint {
    clientSetMutex := self.clientSetMutex
    clientSetMutex.RLock()
    count := uint(len(self.clientSet))
    clientSetMutex.RUnlock()

    return count
}

func (self *Server) RegisterCommand(
        id uint,
        makeFunc func() proto.Message,
        handleFunc func (*proxy.NetProxy, proto.Message)) {
    self.parser.Register(id, makeFunc, handleFunc)
}

func (self *Server) Unregister(id uint) {
    self.parser.Unregister(id)
}

