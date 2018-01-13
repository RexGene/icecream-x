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

const (
	BUFFER_SIZE = 65536
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

type NetProxy struct {
	conn        net.Conn
	buffer      *DataBuffer
	isRunning   bool
	parser      IParser
	netProtocol net_protocol.INetProtocol
	customData  interface{}
	recvicer    ICloseNotifyRecvicer
	bufferMaker *DataBufferMaker
}

func (self *NetProxy) Start() {
	if self.isRunning {
		return
	}

	if self.bufferMaker == nil {
		self.bufferMaker = NewDataBufferMaker(BUFFER_SIZE)
	}

	self.isRunning = true
	go self.read_execute()
}

func (self *NetProxy) Send(cmdId uint, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	newBuffer := self.bufferMaker.GetBuffer()
	defer self.bufferMaker.PutBuffer(newBuffer)

	sendBuffer := newBuffer[:len(data)+HEADER_SIZE]
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
	newBuffer := self.bufferMaker.GetBuffer()
	defer self.bufferMaker.PutBuffer(newBuffer)

	sendBuffer := newBuffer[:len(data)+HEADER_SIZE]
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
			buffer = NewDataBufferByData(self.bufferMaker.GetBuffer())
		}

		size, err := self.netProtocol.Read(buffer.GetDataTail())
		if err != nil {
			if err != io.EOF {
				log.Println("[!]", err)
			}

			self.bufferMaker.PutBuffer(buffer.GetOriginData())
			self.Stop()
			continue
		}

		err = buffer.WriteSize(size)
		if err != nil {
			log.Println("[!]", err)
			self.bufferMaker.PutBuffer(buffer.GetOriginData())
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
			self.bufferMaker.PutBuffer(buffer.GetOriginData())
			continue
		}

		if dataLen < int(header.Len) {
			self.buffer = buffer
			continue
		}

		err = buffer.ReadSize(int(header.Len))
		if err != nil {
			log.Println("[!] header.Len invalid:", header.Len)
			self.bufferMaker.PutBuffer(buffer.GetOriginData())
			continue
		}

		dataLen = buffer.GetDataLen()
		if dataLen == 0 {
			executeData := buffer.GetReadData()

			func() {
				defer self.bufferMaker.PutBuffer(buffer.GetOriginData())

				err = self.read_parseAndHandle(uint(header.CmdId), executeData[HEADER_SIZE:])
				if err != nil {
					log.Println("[!]", err, header.CmdId)
				}
			}()
		} else {
			executeData := buffer.GetReadData()
			SurplusDate := buffer.GetData()

			oldBuffer := buffer
			buffer = NewDataBufferAndCopyData(self.bufferMaker.GetBuffer(), SurplusDate)

			func() {
				defer self.bufferMaker.PutBuffer(oldBuffer.GetOriginData())

				err = self.read_parseAndHandle(uint(header.CmdId), executeData[HEADER_SIZE:])
				if err != nil {
					log.Println("[!]", err, header.CmdId)
				}
			}()
			goto L
		}
	}

	if self.buffer != nil {
		self.bufferMaker.PutBuffer(self.buffer.GetOriginData())
		self.buffer = nil
	}
}

func (self *NetProxy) GetRemoteAddr() string {
	return self.netProtocol.GetRemoteAddr()
}

func (self *NetProxy) GetLocalAddr() string {
	return self.netProtocol.GetLocalAddr()
}
