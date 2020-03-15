package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/safe-k/gonnect"
)

type (
	playerServerStorage interface {
		GetActiveMatch(p gonnect.Player) (*gonnect.Match, error)
		SavePlayer(p *gonnect.Player) error
	}

	playerServer struct {
		ConnectionUpgrader
		Storage playerServerStorage
	}
)

func (s *playerServer) Serve(addr string) {
	r := chi.NewRouter()
	r.Get("/player/match", s.getPlayerMatch)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	log.Println("Listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func (s *playerServer) getPlayerMatch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ws, err := s.Upgrade(w, r)
	if err != nil {
		log.Println("Websocket connection upgrade error:", err)
		return
	}

	defer ws.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Cancel signal received, aborting websocket ping")
				return
			default:
				if err := ws.Ping(); err != nil {
					log.Println("Websocket disconnected prematurely", err)
					cancel()
				}

				<-time.After(5 * time.Second)
			}
		}
	}()

	p := &gonnect.Player{}
	if err := ws.ReceiveJSON(p); err != nil {
		log.Println("Websocket read error:", err)
		return
	}

	matchChan := make(chan *gonnect.Match, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Cancel signal received, aborting match check")
				return
			default:
				m, err := s.Storage.GetActiveMatch(*p)
				if err != nil {
					if !strings.Contains(err.Error(), "no rows") {
						log.Println("Could not fetch match data", err)
						cancel()
						return
					}

					// Interval can be customisable
					<-time.After(2 * time.Second)
					continue
				}

				matchChan <- m
				return
			}
		}
	}()

	p.State = gonnect.PlayerSearching
	if err := s.Storage.SavePlayer(p); err != nil {
		log.Println("Could not update player", err)
		cancel()
	}

	select {
	case <-ctx.Done():
		log.Println("Cancel signal received, aborting request")

		p.State = gonnect.PlayerUnavailable
		if err := s.Storage.SavePlayer(p); err != nil {
			log.Println("Could not update player", err)
		}
	case m := <-matchChan:
		log.Println("Found match for player:", m.ID)

		if err := ws.SendJSON(m); err != nil {
			log.Println("Websocket write error:", err)
		}

		return
	}
}

func PlayerServer(upgrader ConnectionUpgrader, storage playerServerStorage) *playerServer {
	return &playerServer{
		ConnectionUpgrader: upgrader,
		Storage:            storage,
	}
}
