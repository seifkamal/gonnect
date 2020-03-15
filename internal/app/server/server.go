package server

import (
	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/app/server/match"
	"github.com/safe-k/gonnect/internal/app/server/player"
)

func ServePlayer(addr string) {
	storage := internal.Storage()
	defer storage.Close()

	s := &player.Server{Storage: storage}

	s.Serve(addr)
}

func ServeMatch(addr string, auth match.Authenticator) {
	storage := internal.Storage()
	defer storage.Close()

	s := &match.Server{
		Authenticator: auth,
		Storage:       storage,
	}

	s.Serve(addr)
}
