package server

import (
	"github.com/safe-k/gonnect/internal"
)

func ServePlayer(addr string) {
	storage := internal.Storage()
	defer storage.Close()

	s := &playerServer{Storage: storage}

	s.Serve(addr)
}

func ServeMatch(addr string, auth Authenticator) {
	storage := internal.Storage()
	defer storage.Close()

	s := &matchmakingServer{
		Authenticator: auth,
		Storage:       storage,
	}

	s.Serve(addr)
}
