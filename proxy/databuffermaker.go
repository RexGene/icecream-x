package proxy

import (
	"sync"
)

type DataBufferMaker struct {
	pools []sync.Pool
}

func NewDataBufferMaker(size uint) *DataBufferMaker {
	object := &DataBufferMaker{}

	object.pools = make([]sync.Pool, len(_maskTable))

	for i, _ := range object.pools {
		pool := &object.pools[i]
		onNew := func() interface{} {
			return make([]byte, _maskTable[i])
		}
		pool.New = onNew
	}

	return object
}

func (self *DataBufferMaker) GetBuffer(size uint) []byte {
	idx, _ := _getMaskIndex(size)
	return self.pools[idx].Get().([]byte)
}

func (self *DataBufferMaker) PutBuffer(data []byte) {
	idx, _ := _getMaskIndex(uint(len(data)))
	self.pools[idx].Put(data)
}

var _maskTable = []uint16{
	0x1F, 0x3F, 0x7F, 0xFF,
	0x1FF, 0x3FF, 0x7FF, 0xFFF,
	0x1FFF, 0x3FFF, 0x7FFF, 0xFFFF,
}

func _getMaskIndex(v uint) (int, bool) {
	for i, m := range _maskTable {
		if uint16(v) <= m {
			return i, true
		}
	}

	return 0, false
}
