package proxy

import (
    "net"
)

type ClientProxy struct {
    NetProxy
    customData interface{}
}

func NewClientProxy(conn net.Conn, recvicer ICloseNotifyRecvicer, parser IParser) *ClientProxy {
    proxy :=  &ClientProxy {
    }

    proxy.isRunning = false
    proxy.conn = conn
    proxy.parser = parser
    proxy.recvicer = recvicer

    return proxy
}

func (self *ClientProxy) SetCustomData(data interface{}) {
    self.customData = data
}

func (self *ClientProxy) GetCustomData() interface{} {
    return self.customData
}
