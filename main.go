package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/pkg/player/connection"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .server file found")
	}
}

func main() {
	dbUrl := os.Getenv("DB")
	if dbUrl == "" {
		log.Fatalln("DB environment variable not set")
	}

	store, err := internal.ConnectStorage(dbUrl)
	if err != nil {
		log.Fatalln("Could not connect to DB")
	}

	defer store.Disconnect()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Could not complete request:", err)
			return
		}
	}()

	http.Handle("/player/connect", connection.Handler(&store))

	const port = ":5000"
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
