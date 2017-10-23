package net_protocol

import (
	"net"
)

type Tcp struct {
	conn net.Conn
}

func (self *Tcp) Read(data []byte) (int, error) {
	return self.conn.Read(data)
}

func (self *Tcp) Write(data []byte) (int, error) {
	return self.conn.Write(data)
}
