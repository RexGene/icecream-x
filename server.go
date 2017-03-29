package icecreamx

import (
    "net"
    "log"
    "sync"
    "time"
    "github.com/RexGene/icecreamx/proxy"
    "github.com/RexGene/icecreamx/parser"
    "github.com/golang/protobuf/proto"
)

const (
    DEFUALT_CHANNEL_SIZE = 32
)

type Server struct {
    sync.RWMutex
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
    self.isRunning = false
    self.listener.Close()

    close(self.chRemoveClient)

    self.RLock()
    for clientProxy, _ := range self.clientSet {
        clientProxy.(*proxy.ClientProxy).Stop()
    }
    self.RUnlock()

    self.waitGroup.Wait()
}

func (self *Server) PushCloseNotify(v interface{}) {
    if self.isRunning {
        self.chRemoveClient <- v
    }
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

    for self.isRunning {
        conn, err := listener.Accept()
        if err == nil {
            clientProxy := proxy.NewClientProxy(conn, self, self.parser)
            clientProxy.Start()

            self.Lock()
            self.clientSet[clientProxy] = true
            self.Unlock()
        } else {
            log.Println("[!]", err)
        }
    }
}

func (self *Server) eventloop_execute() {
    self.waitGroup.Add(1)
    defer self.waitGroup.Done()

    for self.isRunning {
        select {
            case clientProxy, isOK := <-self.chRemoveClient:
                if isOK {
                    self.Lock()
                    delete(self.clientSet, clientProxy)
                    self.Unlock()
                }
            case <- time.After(time.Second):
        }
    }
}

func (self *Server) init(addr string) error {
    self.addr = addr
    return nil
}

func (self *Server) getConnectionCount() uint {
    self.RLock()
    count := uint(len(self.clientSet))
    self.RUnlock()

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

