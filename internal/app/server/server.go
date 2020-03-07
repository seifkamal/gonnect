package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/safe-k/gonnect/internal/pkg/database"
)

type server struct {
	db *sqlx.DB
}

func Serve() {
	DB := database.New()
	defer DB.Close()

	s := &server{DB}

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
