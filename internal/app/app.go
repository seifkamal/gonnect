package app

import (
	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/app/matchmaking"
	"github.com/safe-k/gonnect/internal/app/server"
	"github.com/safe-k/gonnect/internal/app/server/websocket"
)

func ServePlayer(addr string) {
	storage := internal.Storage()
	defer storage.Close()

	server.PlayerServer(websocket.ConnectionUpgrader(), storage).Serve(addr)
}

func ServeMatchmaking(addr string, auth server.Authenticator) {
	storage := internal.Storage()
	defer storage.Close()

	server.MatchmakingServer(auth, storage).Serve(addr)
}

func MatchPlayers(batch int) {
	storage := internal.Storage()
	defer storage.Close()

	matchmaking.Worker(storage).Work(batch)
}
