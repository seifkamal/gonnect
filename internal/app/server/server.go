package server

import (
	"log"
	"net/http"
)

type Handler interface {
	Router() http.Handler
}

func Serve(h Handler) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	const port = ":5000"
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, h.Router()))
}
