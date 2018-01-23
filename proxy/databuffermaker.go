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

		idx := i
		onNew := func() interface{} {
			return make([]byte, _maskTable[idx])
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

var _maskTable = []uint{
	0x2, 0x10,
	0x20, 0x40, 0x80, 0x100,
	0x200, 0x400, 0x800, 0x1000,
	0x2000, 0x4000, 0x8000, 0x10000,
}

func _getMaskIndex(v uint) (int, bool) {
	for i, m := range _maskTable {
		if v <= m {
			return i, true
		}
	}

	return 0, false
}
