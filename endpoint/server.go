package endpoint

import (
    "net"
    "log"
)

const BUFFER_SIZE = 65535

type Server struct {
	conn        net.Conn
    isRunning   bool
    buffer      *DataBuffer
}


func NewServer(addr string) (*Server, error) {
    server := &Server {
        isRunning : false,
    }

    err := server.init(addr)
    if err != nil {
        return nil, err
    }

    return server, nil
}

func (self *Server) Start() {
    self.isRunning = true
}

func (self *Server) Stop() {
    self.isRunning = false
}

func (self *Server) init(addr string) error {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return err
    }

    self.conn = conn

    return nil
}

func (self *Server) listen_makeBuffer(size uint) []byte{
    return make([]byte, size)
}

func (self *Server) listen_freeBuffer(buffer []byte) {

}

func (self *Server) listen_execute() {
	for self.isRunning {
        var buffer *DataBuffer
        if self.buffer == nil {
            buffer = self.buffer
            self.buffer = nil
        } else {
            buffer = NewDataBufferByData(self.listen_makeBuffer(BUFFER_SIZE))
        }

		size, err := self.conn.Read(buffer.GetDataTail())
		if err == nil {
            err = buffer.WriteSize(size)
            if err != nil {
                log.Println("[!]", err)
                continue
            }

            if buffer.GetDataLen() >= HEADER_SIZE {
                header := CheckSumAndGetHeader(buffer.GetDataHead())
                if header != nil {
                    err = buffer.ReadSize(int(header.Len))
                    if err != nil {
                        log.Println("[!] header.Len invalid:", header.Len)
                        buffer.Reset()
                    }

                    dataLen := buffer.GetDataLen()
                    if dataLen == 0 {
                        // [TODO] buffer execute
                        // buffer.GetReadBuffer() 
                    } else {
                        // [TODO] buffer copy send send
                        // buffer.GetReadBuffer() 
                        self.buffer = buffer
                    }
                } else {
                    buffer.Reset()
                }
            } else {
                self.buffer = buffer
            }
		} else {
			log.Println("[!]", err)
		}
	}
}
