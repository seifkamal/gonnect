package server

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type webSocket struct {
	conn *websocket.Conn
}

func WebSocket(w http.ResponseWriter, r *http.Request) (*webSocket, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if c != nil {
		log.Printf("%s connected\n", c.RemoteAddr())
	}

	return &webSocket{c}, err
}

func (s *webSocket) Address() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *webSocket) SendJSON(i interface{}) error {
	return s.conn.WriteJSON(i)
}

func (s *webSocket) ReceiveJSON(i interface{}) error {
	return s.conn.ReadJSON(i)
}

func (s *webSocket) Ping() error {
	return s.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (s *webSocket) Close() error {
	log.Printf("Closing connection to %s\n", s.Address())

	return s.conn.Close()
}
