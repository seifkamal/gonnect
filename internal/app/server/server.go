package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/pop"
)

type server struct {
	db *pop.Connection
}

func Serve() {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalln("Could not connect to DB", err)
	}
	defer db.Close()

	s := &server{db}

	r := chi.NewRouter()
	r.Get("/player/connect", s.handlePlayerConnect())
	r.Get("/match/all", s.handleGetReadyMatch())

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
