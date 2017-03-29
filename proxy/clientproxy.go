package proxy

import (
    "net"
)

type ICloseNotifyRecvicer interface {
    PushCloseNotify(interface{})
}

type ClientProxy struct {
    NetProxy
    recvicer ICloseNotifyRecvicer
}

func NewClientProxy(conn net.Conn, recvicer ICloseNotifyRecvicer, parser IParser) *ClientProxy {
    proxy :=  &ClientProxy {
        recvicer : recvicer,
    }

    proxy.isRunning = false
    proxy.conn = conn
    proxy.parser = parser

    return proxy
}

func (self *ClientProxy) Stop() {
    self.isRunning = false
    self.conn.Close()

    self.recvicer.PushCloseNotify(self)
}

