package net_protocol

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

const (
	READ_TIMEOUT_S   = 10
	WRITE_TIMEOUT_S  = 10
	MAX_HEADER_BYTES = 1 << 20
)

type WebSocket struct {
	conn *websocket.Conn
}

func NewWebSocket(conn *websocket.Conn) *WebSocket {
	return &WebSocket{
		conn: conn,
	}
}

func (self *WebSocket) Read(data []byte) (int, error) {
	return self.conn.Read(data)
}

func (self *WebSocket) Write(data []byte) (int, error) {
	return self.conn.Write(data)
}

func (self *WebSocket) Close() error {
	return self.conn.Close()
}

type WebSocketListener struct {
	server *http.Server
}

func NewWebSocketListener() *WebSocketListener {
	return &WebSocketListener{}
}

func (self *WebSocketListener) Listen(addr string, handler func(INetProtocol)) error {
	onHandle := func(conn *websocket.Conn) {
		handler(NewWebSocket(conn))
	}

	s := &http.Server{
		Addr:           addr,
		Handler:        websocket.Handler(onHandle),
		ReadTimeout:    READ_TIMEOUT_S * time.Second,
		WriteTimeout:   WRITE_TIMEOUT_S * time.Second,
		MaxHeaderBytes: MAX_HEADER_BYTES,
	}
	self.server = s
	return s.ListenAndServe()
}

func (self *WebSocketListener) Close() error {
	return self.server.Close()
}