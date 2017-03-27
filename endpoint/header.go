package endpoint

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
)

func CheckSumAndGetHeader(buffer []byte) *Header {
	header := (*Header)(unsafe.Pointer(&buffer[0]))
	sum := header.Sum

	header.Sum = 0

	buf := buffer[:header.Len]
	sumValue := byte(0)
	for _, v := range buf {
		sumValue ^= v
	}

	if sum != sumValue {
		log.Println("[!] check sum invaild: sum", sum, "sumValue:", sumValue)
		return nil
	}

	return header
}
