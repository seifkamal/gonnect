package server

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/safe-k/gonnect/internal/domain"
)

type courier struct {
	Player domain.Player `json:"player"`
	Match  domain.Match  `json:"match"`
}

type socket struct {
	conn *websocket.Conn
}

func OpenSocket(w http.ResponseWriter, r *http.Request) (*socket, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if c != nil {
		log.Printf("%s connected\n", c.RemoteAddr())
	}

	return &socket{c}, err
}

func (s *socket) Address() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *socket) Send(c *courier) error {
	return s.conn.WriteJSON(c)
}

func (s *socket) Receive(c *courier) (*courier, error) {
	if c == nil {
		c = &courier{}
	}

	if err := s.conn.ReadJSON(c); err != nil {
		return c, err
	}

	return c, nil
}

func (s *socket) Ping() error {
	return s.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (s *socket) Close() error {
	log.Printf("Closing connection to %s\n", s.Address())

	return s.conn.Close()
}
