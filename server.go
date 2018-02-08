package icecreamx

import (
	"github.com/RexGene/icecreamx/net_protocol"
	"github.com/RexGene/icecreamx/parser"
	"github.com/RexGene/icecreamx/proxy"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
)

const (
	DEFUALT_CHANNEL_SIZE = 10240
)

const (
	BUFFER_SIZE = 65536
)

type Server struct {
	clientSetMutex  sync.RWMutex
	runningMutex    sync.RWMutex
	waitGroup       sync.WaitGroup
	isRunning       bool
	addr            string
	listener        net_protocol.IListener
	parser          *parser.PbParser
	clientSet       map[interface{}]struct{}
	chRemoveClient  chan interface{}
	dataBufferMaker *proxy.DataBufferMaker
}

func NewServer(addr string,
	listener net_protocol.IListener) (*Server, error) {
	server := &Server{
		isRunning: false,
		listener:  listener,
		parser:    parser.NewPbParser(),
	}

	err := server.init(addr)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (self *Server) Start() {
	if !self.isRunning {
		self.isRunning = true
		self.clientSet = make(map[interface{}]struct{})
		self.chRemoveClient = make(chan interface{}, DEFUALT_CHANNEL_SIZE)
		self.dataBufferMaker = proxy.NewDataBufferMaker(BUFFER_SIZE)

		go self.listen_execute()
		go self.eventloop_execute()
	}
}

func (self *Server) StartAndWait() {
	if !self.isRunning {
		self.isRunning = true
		self.clientSet = make(map[interface{}]struct{})
		self.chRemoveClient = make(chan interface{}, DEFUALT_CHANNEL_SIZE)
		self.dataBufferMaker = proxy.NewDataBufferMaker(BUFFER_SIZE)

		go self.listen_execute()
		self.eventloop_execute()
	}
}

func (self *Server) Stop() {
	runningMutex := &self.runningMutex
	runningMutex.Lock()
	self.isRunning = false
	runningMutex.Unlock()

	if self.listener != nil {
		self.listener.Close()
	}

	close(self.chRemoveClient)

	clientSetMutex := &self.clientSetMutex
	clientSetMutex.RLock()
	for clientProxy, _ := range self.clientSet {
		clientProxy.(*proxy.ClientProxy).Stop()
	}
	clientSetMutex.RUnlock()

	self.waitGroup.Wait()
}

func (self *Server) PushCloseNotify(v interface{}) {
	runningMutex := &self.runningMutex
	runningMutex.RLock()
	if self.isRunning {
		self.chRemoveClient <- v
	}
	runningMutex.RUnlock()
}

func (self *Server) listen_execute() {
	self.waitGroup.Add(1)
	defer self.waitGroup.Done()

	log.Println("[*] listening...")
	onNewConn := func(np net_protocol.INetProtocol) {
		clientProxy := proxy.NewClientProxy(np, self, self.parser)

		clientSetMutex := &self.clientSetMutex
		clientSetMutex.Lock()
		self.clientSet[clientProxy] = struct{}{}
		clientSetMutex.Unlock()

		clientProxy.Setup(self.dataBufferMaker)
		np.Start(clientProxy)
	}

	err := self.listener.Listen(self.addr, onNewConn)
	if err != nil {
		log.Fatalln("[-]", err)
	}
}

func (self *Server) eventloop_execute() {
	self.waitGroup.Add(1)
	defer self.waitGroup.Done()

	for self.isRunning {
		select {
		case clientProxy, isOK := <-self.chRemoveClient:
			if isOK {
				clientSetMutex := &self.clientSetMutex
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
	clientSetMutex := &self.clientSetMutex
	clientSetMutex.RLock()
	count := uint(len(self.clientSet))
	clientSetMutex.RUnlock()

	return count
}

func (self *Server) RegisterCommand(
	id uint,
	makeFunc func() proto.Message,
	handleFunc func(*proxy.NetProxy, proto.Message)) {
	self.parser.Register(id, makeFunc, handleFunc)
}

func (self *Server) UnregisterCommand(id uint) {
	self.parser.Unregister(id)
}
