package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/safe-k/gonnect/internal/pkg/match"
	"github.com/safe-k/gonnect/internal/pkg/player"
)

func (s *server) handlePlayerConnect() http.HandlerFunc {
	var (
		once sync.Once
		pr   player.Repository
		mr   match.Repository
	)

	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			pr = player.Repository{DB: s.db}
			mr = match.Repository{DB: s.db}
		})

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

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

		alias := cour.Player.Alias
		matchChan := make(chan int, 1)

		go func() {
			for {
				select {
				case <-ctx.Done():
					log.Println("Cancel signal received, aborting match check")
					return
				default:
					m, err := mr.FindByPlayerAlias(alias)
					if err != nil {
						switch err {
						case sql.ErrNoRows:
							// Interval can be customisable
							<-time.After(2 * time.Second)
							continue
						default:
							log.Println("Could not fetch match data", err)
							cancel()
						}
					}

					matchChan <- int(m.ID)
					return
				}
			}
		}()

		go func() {
			select {
			case <-ctx.Done():
				log.Println("Cancel signal received, aborting player update")
				return
			default:
				_, err := pr.FindByAlias(alias)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						log.Println("Creating new player:", alias)

						_, err = pr.New(alias, player.Searching)
						if err != nil {
							log.Println("Could not create player", err)
							cancel()
						}
					default:
						log.Println("Could not fetch player data", err)
						cancel()
					}
				}
			}
		}()

		select {
		case <-ctx.Done():
			log.Println("Cancel signal received, aborting request")
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

func (s *server) handleGetReadyMatch() http.HandlerFunc {
	var (
		once sync.Once
		mr   match.Repository
	)

	return func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			mr = match.Repository{DB: s.db}
		})

		st := r.URL.Query().Get("state")
		if st == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("'state' query parameter required"))
			if err != nil {
				log.Println("Could not send error response", err)
				return
			}

			return
		}

		mm, err := mr.Find(st)
		if err != nil {
			log.Println("Could not find ready matches", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		mmBytes, err := json.Marshal(mm)
		if err != nil {
			log.Println("Could not JSON encode match data", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = w.Write(mmBytes)
		if err != nil {
			log.Println("Could not send match data response", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
