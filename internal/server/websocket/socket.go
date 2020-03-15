package websocket

import (
	"log"
	"net"

	ws "github.com/gorilla/websocket"
)

type socket struct {
	conn *ws.Conn
}

func (s *socket) Address() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *socket) SendJSON(i interface{}) error {
	return s.conn.WriteJSON(i)
}

func (s *socket) ReceiveJSON(i interface{}) error {
	return s.conn.ReadJSON(i)
}

func (s *socket) Ping() error {
	return s.conn.WriteMessage(ws.PingMessage, []byte{})
}

func (s *socket) Close() error {
	log.Printf("Closing connection to %s\n", s.Address())

	return s.conn.Close()
}
