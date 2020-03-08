package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/safe-k/gonnect/internal"
)

type server internal.Server

func Serve() {
	db := internal.DB()
	defer db.Close()

	s := &server{DB: db}

	r := chi.NewRouter()
	r.Get("/player/connect", s.connectPlayer())
	r.Route("/match", func(r chi.Router) {
		r.Get("/all", s.getAllMatches())
		r.Route("/{matchId}", func(r chi.Router) {
			r.Get("/", s.getMatch())
			r.Post("/end", s.endMatch())
		})
	})

	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	const port = ":5000"
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
