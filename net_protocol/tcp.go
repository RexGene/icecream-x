package net_protocol

import (
	"log"
	"net"
)

type Tcp struct {
	conn      net.Conn
	listenner net.Listener
	isRunning bool
}

func NewTcp(conn net.Conn) *Tcp {
	object := &Tcp{
		conn:      conn,
		isRunning: false,
	}

	return object
}

func (self *Tcp) Read(data []byte) (int, error) {
	return self.conn.Read(data)
}

func (self *Tcp) Write(data []byte) (int, error) {
	return self.conn.Write(data)
}

func (self *Tcp) GetRemoteAddr() string {
	return self.conn.RemoteAddr().String()
}

func (self *Tcp) GetLocalAddr() string {
	return self.conn.LocalAddr().String()
}

func (self *Tcp) Close() error {
	self.isRunning = false
	return self.conn.Close()
}

type TcpListener struct {
	isRunning bool
	listener  net.Listener
}

func NewTcpListener() *TcpListener {
	return &TcpListener{
		isRunning: false,
	}
}

func (self *TcpListener) Close() error {
	self.isRunning = false
	return self.listener.Close()
}

func (self *TcpListener) Listen(addr string, handler func(np INetProtocol)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("[-]", err)
		return err
	}
	defer listener.Close()

	self.listener = listener

	self.isRunning = true
	for self.isRunning {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[!]", err)
			continue
		}

		handler(NewTcp(conn))
	}

	return nil
}
