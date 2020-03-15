package server

import (
	"log"
	"net"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type websocket struct {
	conn *ws.Conn
}

func Websocket(w http.ResponseWriter, r *http.Request) (*websocket, error) {
	upgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if conn != nil {
		log.Printf("%s connected\n", conn.RemoteAddr())
	}

	return &websocket{conn}, err
}

func (s *websocket) Address() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *websocket) SendJSON(i interface{}) error {
	return s.conn.WriteJSON(i)
}

func (s *websocket) ReceiveJSON(i interface{}) error {
	return s.conn.ReadJSON(i)
}

func (s *websocket) Ping() error {
	return s.conn.WriteMessage(ws.PingMessage, []byte{})
}

func (s *websocket) Close() error {
	log.Printf("Closing connection to %s\n", s.Address())

	return s.conn.Close()
}
