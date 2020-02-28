package server

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/safe-k/gonnect/internal/pkg"
)

func (s *server) handlePlayerConnect() http.HandlerFunc {
	var (
		once sync.Once
		pr   pkg.PlayerRepository
		mr   pkg.MatchRepository
	)

	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			pr = pkg.PlayerRepository{DB: s.db}
			mr = pkg.MatchRepository{DB: s.db}
		})

		ws, err := OpenSocket(w, r)
		if err != nil {
			log.Println("WebSocket connection upgrade error:", err)
			return
		}

		defer ws.Close()

		cour, err := ws.Receive(nil)
		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}

		dcChan := make(chan error, 1)
		go func() {
			for {
				if err := ws.Ping(); err != nil {
					dcChan <- err
					return
				}

				<-time.After(5 * time.Second)
			}
		}()

		matchChan := make(chan int, 1)
		errChan := make(chan error, 1)
		go func() {
			alias := cour.Player.Alias

			p, err := pr.FindByAlias(alias)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					log.Println("Creating new player:", alias)

					p, err = pr.New(alias)
					if err != nil {
						log.Println("Could not create player:", err)
						errChan <- err
						return
					}
				default:
					log.Println("Error finding player:", err)
					errChan <- err
					return
				}
			} else if p.State != pkg.PlayerOnline {
				p.State = pkg.PlayerOnline
				if err := pr.Save(p); err != nil {
					log.Println("Could not update player state:", err)
					errChan <- err
					return
				}
			}

			log.Println("Checking player state")
			for {
				m, err := mr.FindByPlayerAlias(alias)
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
		}()

		select {
		case err := <-dcChan:
			log.Panicln("User disconnected prematurely:", err)
		case err := <-errChan:
			log.Println("Could not find match for player:", err)
		case matchID := <-matchChan:
			log.Println("Found match for player:", matchID)

			cour.Match.ID = matchID
			if err := ws.Send(cour); err != nil {
				log.Println("WebSocket write error:", err)
			}

			return
		}
	}
}
