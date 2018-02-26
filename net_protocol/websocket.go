package net_protocol

import (
	// websocket "github.com/RexGene/websocket-go"
	// "golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	READ_TIMEOUT_S   = 10
	WRITE_TIMEOUT_S  = 10
	MAX_HEADER_BYTES = 1 << 20
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

type Handler func(conn *websocket.Conn)

func (self Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[-] Upgrade:", err)
		return

	}

	defer conn.Close()
	self(conn)
}

type WebSocket struct {
	conn       *websocket.Conn
	lastReader io.Reader
}

func NewWebSocket(conn *websocket.Conn) *WebSocket {
	return &WebSocket{
		conn: conn,
	}
}

func (self *WebSocket) Read(data []byte) (int, error) {
	if self.lastReader != nil {
		ret, err := self.lastReader.Read(data)
		if err == nil {
			return ret, err
		}

		if err != io.EOF {
			log.Println("[!]", err)
		}
		self.lastReader = nil
	}

	_, r, err := self.conn.NextReader()
	if err != nil {
		return 0, err
	}

	self.lastReader = r
	return r.Read(data)
}

func (self *WebSocket) Write(data []byte) (int, error) {
	w, err := self.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}

	defer w.Close()
	return w.Write(data)
}

func (self *WebSocket) Close() error {
	return self.conn.Close()
}

func (self *WebSocket) GetRemoteAddr() string {
	return self.conn.RemoteAddr().String()
}

func (self *WebSocket) GetLocalAddr() string {
	return self.conn.LocalAddr().String()
}

func (self *WebSocket) Start(selector Selector) {
	selector.StartAndWait()
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
		Handler:        Handler(onHandle),
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

type WebSocketSecureListener struct {
	server      *http.Server
	crtFilePath string
	keyFilePath string
}

func NewWebSocketSecureListener(crtFilePath, keyFilePath string) *WebSocketSecureListener {
	return &WebSocketSecureListener{
		crtFilePath: crtFilePath,
		keyFilePath: keyFilePath,
	}
}

func (self *WebSocketSecureListener) Listen(addr string, handler func(INetProtocol)) error {
	onHandle := func(conn *websocket.Conn) {
		handler(NewWebSocket(conn))
	}

	s := &http.Server{
		Addr:           addr,
		Handler:        Handler(onHandle),
		ReadTimeout:    READ_TIMEOUT_S * time.Second,
		WriteTimeout:   WRITE_TIMEOUT_S * time.Second,
		MaxHeaderBytes: MAX_HEADER_BYTES,
	}
	self.server = s
	return s.ListenAndServeTLS(self.crtFilePath, self.keyFilePath)
}

func (self *WebSocketSecureListener) Close() error {
	return self.server.Close()
}
