package proxy

import (
	"errors"
)

var (
	ErrParamInvalid   = errors.New("<DataBuffer> param invalid")
	ErrOffsetOverflow = errors.New("<DataBuffer> offset overflow")
)

type DataBuffer struct {
	data        []byte
	writeOffset int
	readOffset  int
}

func NewDataBufferByData(data []byte) *DataBuffer {
	return &DataBuffer{
		data: data,
	}
}

func NewDataBufferAndCopyData(newBuffer []byte, data []byte) *DataBuffer {
	buffer := &DataBuffer{
		data: newBuffer,
	}

	copy(buffer.data, data)
	buffer.writeOffset = len(data)

	return buffer
}

func (self *DataBuffer) GetData() []byte {
	return self.data[self.readOffset:self.writeOffset]
}

func (self *DataBuffer) GetDataHead() []byte {
	return self.data[self.readOffset:]
}

func (self *DataBuffer) GetDataTail() []byte {
	return self.data[self.writeOffset:]
}

func (self *DataBuffer) ReadSize(size int) error {
	if size < 0 {
		return ErrParamInvalid
	}

	value := self.readOffset + size
	if value < self.readOffset || value > self.writeOffset {
		return ErrOffsetOverflow
	}

	self.readOffset = value
	return nil
}

func (self *DataBuffer) WriteSize(size int) error {
	if size < 0 {
		return ErrParamInvalid
	}

	value := self.writeOffset + size
	if value < self.writeOffset || value > len(self.data) {
		return ErrOffsetOverflow
	}

	self.writeOffset = value
	return nil
}

func (self *DataBuffer) GetReadOffset() int {
	return self.readOffset
}

func (self *DataBuffer) GetWriteOffset() int {
	return self.writeOffset
}

func (self *DataBuffer) GetReadData() []byte {
	return self.data[:self.readOffset]
}

func (self *DataBuffer) GetDataLen() int {
	return self.writeOffset - self.readOffset
}

func (self *DataBuffer) Reset() {
	self.writeOffset = 0
	self.readOffset = 0
}

func (self *DataBuffer) GetOriginData() []byte {
	return self.data
}
