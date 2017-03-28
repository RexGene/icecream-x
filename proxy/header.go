package proxy

import (
    "unsafe"
    "log"
)

type Header struct {
    Version  uint8
	Sum      uint8
	Len      uint16
    CmdId    uint16
}

const (
    HEADER_SIZE = 6
    VERSION     = 1
)

func CheckSumAndGetHeader(buffer []byte) *Header {
	header := (*Header)(unsafe.Pointer(&buffer[0]))

	buf := buffer[:header.Len]
	sumValue := byte(0)
	for _, v := range buf {
		sumValue ^= v
	}

	if sumValue != 0 {
		log.Println("[!] check sum invaild sumValue:", sumValue)
		return nil
	}

	return header
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
