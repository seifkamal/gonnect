package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/seifkamal/gonnect"
)

type (
	matchmakingServerStorage interface {
		GetMatchesByState(state string) (*gonnect.Matches, error)
		GetMatchById(id int) (*gonnect.Match, error)
		EndMatch(id int) error
	}

	matchmakingServer struct {
		Authenticator
		Storage matchmakingServerStorage
	}
)

func (s *matchmakingServer) Serve(addr string) {
	r := chi.NewRouter()
	r.Route("/match", func(r chi.Router) {
		r.Get("/all", s.getAllMatches)
		r.Route("/{matchId}", func(r chi.Router) {
			r.Get("/", s.getMatch)
			r.Post("/end", s.Authenticate(s.endMatch))
		})
	})

	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	log.Println("Listening on port", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func (s *matchmakingServer) getAllMatches(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("'state' query parameter required"))
		if err != nil {
			log.Println("Could not send error response", err)
			return
		}

		return
	}

	mm, err := s.Storage.GetMatchesByState(state)
	if err != nil {
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

func (s *matchmakingServer) getMatch(w http.ResponseWriter, r *http.Request) {
	mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
	if err != nil {
		log.Println("Could not parse match ID param:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m, err := s.Storage.GetMatchById(mID)
	if err != nil {
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

func (s *matchmakingServer) endMatch(w http.ResponseWriter, r *http.Request) {
	mID, err := strconv.Atoi(chi.URLParam(r, "matchId"))
	if err != nil {
		log.Println("Could not parse match ID param:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Storage.EndMatch(mID)
	if err != nil {
		log.Println("Could not update match:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func MatchmakingServer(auth Authenticator, storage matchmakingServerStorage) *matchmakingServer {
	return &matchmakingServer{
		Authenticator: auth,
		Storage:       storage,
	}
}
