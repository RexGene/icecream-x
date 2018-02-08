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
	"unsafe"
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
	conn         net.Conn
	isRunning    bool
	parser       IParser
	netProtocol  net_protocol.INetProtocol
	customData   interface{}
	recvicer     ICloseNotifyRecvicer
	bufferMaker  *DataBufferMaker
	headerBuffer *DataBuffer
}

func (self *NetProxy) Setup(bufferMaker *DataBufferMaker) {
	self.bufferMaker = bufferMaker
	self.headerBuffer = NewDataBufferByData(self.bufferMaker.GetBuffer(HEADER_SIZE), HEADER_SIZE)
}

func (self *NetProxy) Start() {
	if self.isRunning {
		return
	}

	self.isRunning = true
	go self.read_execute()
}

func (self *NetProxy) StartAndWait() {
	if self.isRunning {
		return
	}

	self.isRunning = true
	self.read_execute()
}

func (self *NetProxy) Send(cmdId uint, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	size := uint(len(data) + HEADER_SIZE)
	newBuffer := self.bufferMaker.GetBuffer(size)
	defer self.bufferMaker.PutBuffer(newBuffer)

	sendBuffer := newBuffer[:size]
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
	size := uint(len(data) + HEADER_SIZE)
	newBuffer := self.bufferMaker.GetBuffer(size)
	defer self.bufferMaker.PutBuffer(newBuffer)

	sendBuffer := newBuffer[:size]
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
		if !self._readHeader() {
			continue
		}

		self._readData()
	}
}

func (self *NetProxy) GetRemoteAddr() string {
	return self.netProtocol.GetRemoteAddr()
}

func (self *NetProxy) GetLocalAddr() string {
	return self.netProtocol.GetLocalAddr()
}

func (self *NetProxy) _readHeader() bool {
	buffer := self.headerBuffer

	for self.isRunning {
		size, err := self.netProtocol.Read(buffer.GetDataTailWithSize())
		if err != nil {
			if err != io.EOF {
				log.Println("[!]", err)
			}

			log.Println("[!] err:", err, "size:", size)
			buffer.Reset()
			self.Stop()
			return false
		}

		err = buffer.WriteSize(size)
		if err != nil {
			log.Println("[!]", err)
			buffer.Reset()
			return false
		}

		dataLen := buffer.GetDataLen()
		if dataLen == HEADER_SIZE {
			return true
		}
	}

	buffer.Reset()
	return false
}

func (self *NetProxy) _getHeader() *Header {
	headerBuffer := self.headerBuffer
	if headerBuffer.GetDataLen() != HEADER_SIZE {
		return nil
	}

	return (*Header)(unsafe.Pointer(&headerBuffer.GetData()[0]))
}

func (self *NetProxy) _readData() bool {
	header := self._getHeader()
	defer self.headerBuffer.Reset()
	if header == nil {
		return false
	}

	dataSize := uint(header.Len) - HEADER_SIZE

	buffer := NewDataBufferByData(self.bufferMaker.GetBuffer(dataSize), int(dataSize))
	defer self.bufferMaker.PutBuffer(buffer.GetOriginData())

	for self.isRunning {
		if dataSize != 0 {
			size, err := self.netProtocol.Read(buffer.GetDataTailWithSize())
			if err != nil {
				if err != io.EOF {
					log.Println("[!]", err)
				}

				log.Println("[!] err:", err, "size:", size)
				self.Stop()
				return false
			}

			buffer.WriteSize(size)
		}

		dataLen := uint(buffer.GetDataLen())
		if dataLen < dataSize {
			continue
		}

		executeData := buffer.GetData()
		if !CheckSum(self.headerBuffer.GetData(), buffer.GetData()) {
			log.Println("[!] check sum error")
			return false
		}

		err := self.read_parseAndHandle(uint(header.CmdId), executeData)
		if err != nil {
			log.Println("[!]", err, header.CmdId)
		}
		break
	}

	return self.isRunning
}
