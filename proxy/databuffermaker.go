package proxy

import (
	"sync"
)

type DataBufferMaker struct {
	pool sync.Pool
}

func NewDataBufferMaker(size uint) *DataBufferMaker {
	object := &DataBufferMaker{}
	onNew := func() interface{} {
		return make([]byte, size)
	}
	object.pool.New = onNew
	return object
}

func (self *DataBufferMaker) GetBuffer() []byte {
	return self.pool.Get().([]byte)
}

func (self *DataBufferMaker) PutBuffer(data []byte) {
	self.pool.Put(data)
}
