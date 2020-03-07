package server

import (
	"log"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.HandleFunc("/player/connect", s.handlePlayerConnect())
	r.HandleFunc("/matches", s.handleGetReadyMatch()).Methods("GET")
	http.Handle("/", r)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	const port = ":5000"
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
