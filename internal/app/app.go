package app

import (
	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/app/matchmaking"
	"github.com/safe-k/gonnect/internal/app/server"
)

func ServePlayer(addr string) {
	storage := internal.Storage()
	defer storage.Close()

	s := &server.PlayerServer{Storage: storage}

	s.Serve(addr)
}

func ServeMatchmaking(addr string, auth server.Authenticator) {
	storage := internal.Storage()
	defer storage.Close()

	s := &server.MatchmakingServer{
		Authenticator: auth,
		Storage:       storage,
	}

	s.Serve(addr)
}

func MatchPlayers(batch int) {
	storage := internal.Storage()
	defer storage.Close()

	w := &matchmaking.Worker{Storage: storage}
	w.Work(batch)
}
