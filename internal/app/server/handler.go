package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/safe-k/gonnect/internal/domain"
)

func (s *server) handlePlayerConnect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
					q := s.db.Q().
						LeftJoin("matches_players", "matches_players.match_id=matches.id").
						LeftJoin("players", "players.id=matches_players.player_id").
						Where("players.alias = ?", alias).
						Where("matches.state <> ?", domain.MatchEnded)

					m := &domain.Match{}
					if err := q.First(m); err != nil {
						if !strings.Contains(err.Error(), "no rows") {
							log.Println("Could not fetch match data", err)
							cancel()
							return
						}

						// Interval can be customisable
						<-time.After(2 * time.Second)
						continue
					}

					matchChan <- m.ID
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
				p := &domain.Player{}
				if err := s.db.Where("alias = ?", alias).First(p); err != nil {
					if !strings.Contains(err.Error(), "no rows") {
						log.Println("Could not fetch player data", err)
						cancel()
						return
					}

					log.Println("Creating new player:", alias)
					p.Alias = alias
				}

				p.State = domain.PlayerSearching
				if err := s.db.Save(p); err != nil {
					log.Println("Could not create player", err)
					cancel()
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
	return func(w http.ResponseWriter, r *http.Request) {
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

		mm := &[]domain.Match{}
		if err := s.db.Where("state = ?", st).All(mm); err != nil {
			log.Println("Could not find ready matches", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mmBytes, err := json.Marshal(mm)
		if err != nil {
			log.Println("Could not JSON encode match data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(mmBytes)
		if err != nil {
			log.Println("Could not send match data response", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleGetMatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
		if err != nil {
			log.Println("Could not parse match ID param:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m := &domain.Match{}
		if err := s.db.Find(m, mID); err != nil {
			if strings.Contains(err.Error(), "no rows") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Println("Could not find match:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		res, err := json.Marshal(m)
		if err != nil {
			log.Println("Could not JSON marshall match data:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			log.Println("Could not send match data:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleEndMatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
		if err != nil {
			log.Println("Could not parse match ID param:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m := &domain.Match{ID: mID, State: domain.MatchEnded}
		if err := s.db.Update(m); err != nil {
			log.Println("Could not update match:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
