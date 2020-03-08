package match

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/safe-k/gonnect/internal/app"
	"github.com/safe-k/gonnect/internal/domain"
)

type Handler app.Actor

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/match", func(r chi.Router) {
		r.Get("/all", h.getAllMatches)
		r.Route("/{matchId}", func(r chi.Router) {
			r.Get("/", h.getMatch)
			r.Post("/end", h.endMatch)
		})
	})
	return r
}

func (h *Handler) getAllMatches(w http.ResponseWriter, r *http.Request) {
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
	if err := h.DB.Where("state = ?", st).All(mm); err != nil {
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

func (h *Handler) getMatch(w http.ResponseWriter, r *http.Request) {
	mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
	if err != nil {
		log.Println("Could not parse match ID param:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m := &domain.Match{}
	if err := h.DB.Find(m, mID); err != nil {
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

func (h *Handler) endMatch(w http.ResponseWriter, r *http.Request) {
	mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
	if err != nil {
		log.Println("Could not parse match ID param:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m := &domain.Match{ID: mID, State: domain.MatchEnded}
	if err := h.DB.Update(m); err != nil {
		log.Println("Could not update match:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
