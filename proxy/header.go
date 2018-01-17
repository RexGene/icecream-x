package proxy

import (
	"unsafe"
)

type Header struct {
	Version uint8
	Sum     uint8
	Len     uint16
	CmdId   uint16
}

const (
	HEADER_SIZE = 6
	VERSION     = 1
)

func CheckSum(headerData []byte, data []byte) bool {
	header := (*Header)(unsafe.Pointer(&headerData[0]))
	if header.Len < uint16(len(headerData)) {
		return false
	}

	buf := data[:(header.Len - uint16(len(headerData)))]
	sumValue := byte(0)

	for _, v := range headerData {
		sumValue ^= v
	}

	for _, v := range buf {
		sumValue ^= v
	}

	if sumValue != 0 {
		return false
	}

	return true
}

func FillHeader(cmdId uint, buffer []byte) {
	header := (*Header)(unsafe.Pointer(&buffer[0]))
	header.Version = VERSION
	header.Len = uint16(len(buffer))
	header.CmdId = uint16(cmdId)
	header.Sum = 0

	sumValue := byte(0)
	for _, v := range buffer {
		sumValue ^= v
	}

	header.Sum = sumValue
}
