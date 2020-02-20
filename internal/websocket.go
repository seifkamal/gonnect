package internal

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type courier struct {
	Player struct{ Alias string `json:"alias"` } `json:"player"`
	Match  struct{ ID int `json:"id"` }          `json:"match"`
}

type Socket interface {
	Address() net.Addr
	Send(c *courier) error
	Receive(c *courier) (*courier, error)
	Ping() error
	Close() error
}

type socket struct {
	conn *websocket.Conn
}

func OpenSocket(w http.ResponseWriter, r *http.Request) (Socket, error) {
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
