package proxy

import (
	"errors"
	"github.com/RexGene/icecreamx/net_protocol"
	"github.com/RexGene/icecreamx/utils"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"runtime/debug"
)

type ICloseNotifyRecvicer interface {
	PushCloseNotify(interface{})
}

type IParser interface {
	ParseAndHandle(proxy *NetProxy, cmdId uint, data []byte) error
}

var (
	ErrCatchException = errors.New("<NetProxy> catch exception")
	ErrRequestTooFast = errors.New("<NetProxy> request too fast")
)

const (
	BUFFER_SIZE = 65536
)

type NetProxy struct {
	conn        net.Conn
	buffer      *DataBuffer
	isRunning   bool
	parser      IParser
	netProtocol net_protocol.INetProtocol
	customData  interface{}
	recvicer    ICloseNotifyRecvicer
}

func (self *NetProxy) Start() {
	if self.isRunning {
		return
	}

	self.isRunning = true
	go self.read_execute()
}

func (self *NetProxy) Send(cmdId uint, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	sendBuffer := make([]byte, len(data)+HEADER_SIZE)
	copy(sendBuffer[HEADER_SIZE:], data)

	FillHeader(cmdId, sendBuffer)

	offset := 0
	dataLen := len(sendBuffer)
	for offset < dataLen {
		wirteSize, err := self.netProtocol.Write(sendBuffer[offset:])
		if err != nil {
			return err
		}

		offset += wirteSize
	}

	return nil
}

func (self *NetProxy) SendData(cmdId uint, data []byte) error {
	sendBuffer := make([]byte, len(data)+HEADER_SIZE)
	copy(sendBuffer[HEADER_SIZE:], data)

	FillHeader(cmdId, sendBuffer)

	offset := 0
	dataLen := len(sendBuffer)
	for offset < dataLen {
		wirteSize, err := self.netProtocol.Write(sendBuffer[offset:])
		if err != nil {
			return err
		}

		offset += wirteSize
	}

	return nil
}

func (self *NetProxy) Stop() {
	if !self.isRunning {
		return
	}

	self.isRunning = false
	self.netProtocol.Close()

	recvicer := self.recvicer
	if recvicer != nil {
		recvicer.PushCloseNotify(self)
	}
}

func (self *NetProxy) SetCustomData(data interface{}) {
	self.customData = data
}

func (self *NetProxy) GetCustomData() interface{} {
	return self.customData
}

func (self *NetProxy) read_parseAndHandle(cmdId uint, data []byte) (err error) {
	defer func() {
		if ex := recover(); ex != nil {
			log.Println("[!] catch exception:", ex)
			debug.PrintStack()
			err = ErrCatchException
		}
	}()

	return self.parser.ParseAndHandle(self, cmdId, data)
}

func (self *NetProxy) read_execute() {
	defer func() {
		utils.PrintRecover(recover())
		self.Stop()
	}()

	for self.isRunning {
		var buffer *DataBuffer
		if self.buffer != nil {
			buffer = self.buffer
			self.buffer = nil
		} else {
			buffer = NewDataBufferByData(make([]byte, BUFFER_SIZE))
		}

		size, err := self.netProtocol.Read(buffer.GetDataTail())
		if err != nil {
			if err != io.EOF {
				log.Println("[!]", err)
			}
			self.Stop()
			continue
		}

		log.Println("[?] ip:", self.GetRemoteAddr())
		err = buffer.WriteSize(size)
		if err != nil {
			log.Println("[!]", err)
			continue
		}

	L:
		dataLen := buffer.GetDataLen()
		if dataLen < HEADER_SIZE {
			self.buffer = buffer
			continue
		}

		header := CheckSumAndGetHeader(buffer.GetDataHead())
		if header == nil {
			log.Println("[!] check sum error")
			buffer.Reset()
			continue
		}

		if dataLen < int(header.Len) {
			self.buffer = buffer
			continue
		}

		err = buffer.ReadSize(int(header.Len))
		if err != nil {
			log.Println("[!] header.Len invalid:", header.Len)
			buffer.Reset()
			continue
		}

		dataLen = buffer.GetDataLen()
		if dataLen == 0 {
			executeData := buffer.GetReadData()
			err = self.read_parseAndHandle(uint(header.CmdId), executeData[HEADER_SIZE:])
			if err != nil {
				log.Println("[!]", err, header.CmdId)
			}
		} else {
			executeData := buffer.GetReadData()
			SurplusDate := buffer.GetData()
			buffer = NewDataBufferAndCopyData(BUFFER_SIZE, SurplusDate)
			err = self.read_parseAndHandle(uint(header.CmdId), executeData[HEADER_SIZE:])
			if err != nil {
				log.Println("[!]", err, header.CmdId)
			}
			goto L
		}
	}
}

func (self *NetProxy) GetRemoteAddr() string {
	return self.netProtocol.GetRemoteAddr()
}

func (self *NetProxy) GetLocalAddr() string {
	return self.netProtocol.GetLocalAddr()
}
