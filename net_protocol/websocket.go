package net_protocol

import (
	"net"
)

type WebSocket struct {
	conn net.Conn
}

func (self *WebSocket) Read(data []byte) (int, error) {
	return 0, nil
}

func (self *WebSocket) Write(data []byte) (int, error) {
	return 0, nil
}
