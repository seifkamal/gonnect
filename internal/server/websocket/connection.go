package websocket

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"

	"github.com/seifkamal/gonnect/internal/server"
)

type connectionUpgrader struct{}

func (c *connectionUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (server.Websocket, error) {
	upgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if conn != nil {
		log.Printf("%s connected\n", conn.RemoteAddr())
	}

	return &socket{conn}, err
}

func ConnectionUpgrader() *connectionUpgrader {
	return &connectionUpgrader{}
}
