package connection

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/pkg/match"
	"github.com/safe-k/gonnect/pkg/player"
)

type handler struct {
	storage *internal.Storage
}

func Handler(s *internal.Storage) http.Handler {
	return &handler{s}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sock, err := internal.OpenSocket(w, r)
	if err != nil {
		log.Println("WebSocket connection upgrade error:", err)
		return
	}

	defer sock.Close()

	cour, err := sock.Receive(nil)
	if err != nil {
		log.Println("WebSocket read error:", err)
		return
	}

	dcChan := make(chan error, 1)
	go h.ping(sock, dcChan)

	matchChan := make(chan int, 1)
	errChan := make(chan error, 1)
	go h.match(cour.Player.Alias, matchChan, errChan)

	select {
	case err := <-dcChan:
		log.Panicln("User disconnected prematurely:", err)
	case err := <-errChan:
		log.Println("Could not find match for player:", err)
	case matchID := <-matchChan:
		log.Println("Found match for player:", matchID)

		cour.Match.ID = matchID
		if err := sock.Send(cour); err != nil {
			log.Println("WebSocket write error:", err)
		}

		return
	}
}

func (h *handler) ping(sock internal.Socket, dcc chan<- error) {
	for {
		if err := sock.Ping(); err != nil {
			dcc <- err
			return
		}

		<-time.After(5 * time.Second)
	}
}

func (h *handler) match(alias string, matchChan chan<- int, errChan chan<- error) {
	pRepo := player.Repository(*h.storage)

	p, err := pRepo.FindByAlias(alias)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			p = player.New(alias)

			log.Println("Creating new player:", p)

			if err := pRepo.Insert(p); err != nil {
				log.Println("Could not persist player:", err)
				errChan <- err
				return
			}
		default:
			log.Println("Error finding player:", err)
			errChan <- err
			return
		}
	} else if p.State != player.Online {
		p.State = player.Online
		if err := pRepo.Update(p); err != nil {
			log.Println("Could not update player state:", err)
			errChan <- err
			return
		}
	}

	log.Println("Checking player state")
	mRepo := match.Repository(*h.storage)
	for {
		m, err := mRepo.FindByPlayerAlias(alias)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				// Interval can be customisable
				<-time.After(2 * time.Second)
				continue
			default:
				errChan <- err
				return
			}
		}

		matchChan <- m.ID
		return
	}
}
