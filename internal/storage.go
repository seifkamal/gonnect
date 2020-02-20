package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Storage
}

type Storage interface {
	FindOne(i interface{}, query string, arg ...interface{}) error
	ExecuteTransaction(query string, arg interface{}) (sql.Result, error)
	Disconnect() error
}

type storage struct {
	DB *sqlx.DB
}

func ConnectStorage(url string) (Storage, error) {
	DB, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, err
	}

	return &storage{DB}, nil
}

func (s *storage) FindOne(i interface{}, query string, args ...interface{}) error {
	return s.DB.Get(i, query, args...)
}

func (s *storage) ExecuteTransaction(query string, arg interface{}) (sql.Result, error) {
	tx := s.DB.MustBegin()

	res, err := tx.NamedExec(query, arg)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *storage) Disconnect() error {
	return s.DB.Close()
}
