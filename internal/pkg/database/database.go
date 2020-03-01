package database

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func New() *sqlx.DB {
	if err := godotenv.Load(); err != nil {
		log.Print("No .server file found")
	}

	dbUrl := os.Getenv("DB")
	if dbUrl == "" {
		log.Fatalln("DB environment variable not set")
	}

	DB, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		log.Fatalln("Could not connect to DB", err)
	}

	return DB
}
