package net_protocol

import ()

type WebSocket struct {
}

func NewWebSocket() *WebSocket {
	return &WebSocket{}
}

func (self *WebSocket) Read(data []byte) (int, error) {
	return 0, nil
}

func (self *WebSocket) Write(data []byte) (int, error) {
	return 0, nil
}

func (self *WebSocket) Close() error {
	return nil
}

func Listen(addr string) (*WebSocket, error) {
	return nil, nil
}
