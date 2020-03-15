package server

import (
	"net"
	"net/http"
)

type (
	Websocket interface {
		Address() net.Addr
		SendJSON(i interface{}) error
		ReceiveJSON(i interface{}) error
		Ping() error
		Close() error
	}

	ConnectionUpgrader interface {
		Upgrade(w http.ResponseWriter, r *http.Request) (Websocket, error)
	}
)
