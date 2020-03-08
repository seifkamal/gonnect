package internal

import (
	"log"

	"github.com/gobuffalo/pop"
)

type Server struct {
	DB *pop.Connection
}

func DB() *pop.Connection {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalln("Could not connect to DB", err)
	}

	return db
}