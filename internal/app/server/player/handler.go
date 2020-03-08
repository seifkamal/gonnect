package player

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/safe-k/gonnect/internal/app"
	"github.com/safe-k/gonnect/internal/domain"
)

type Handler app.Actor

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/player/match", h.getPlayerMatch)
	return r
}

func (h *Handler) getPlayerMatch(w http.ResponseWriter, r *http.Request) {
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
				q := h.DB.Q().
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

	p := &domain.Player{}
	if err := h.DB.Where("alias = ?", alias).First(p); err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			log.Println("Could not fetch player data", err)
			cancel()
			return
		}

		log.Println("Creating new player:", alias)
		p.Alias = alias
	}

	p.State = domain.PlayerSearching
	if err := h.DB.Save(p); err != nil {
		log.Println("Could not update player", err)
		cancel()
	}

	select {
	case <-ctx.Done():
		log.Println("Cancel signal received, aborting request")

		p.State = domain.PlayerUnavailable
		if err := h.DB.Save(p); err != nil {
			log.Println("Could not update player", err)
		}
	case matchID := <-matchChan:
		log.Println("Found match for player:", matchID)

		cour.Match.ID = matchID
		if err := ws.Send(cour); err != nil {
			log.Println("WebSocket write error:", err)
		}

		return
	}
}
