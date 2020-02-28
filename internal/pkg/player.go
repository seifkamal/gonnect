package pkg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	PlayerOffline = "offline"
	PlayerOnline  = "online"
)

type player struct {
	ID    int64  `db:"id"`
	Alias string `db:"alias"`
	State string `db:"state"`
}

type PlayerRepository struct {
	*sqlx.DB
}

func (pr *PlayerRepository) New(alias string) (*player, error) {
	res, err := pr.Exec(`INSERT INTO player (alias, state) VALUES(?, ?)`, alias, PlayerOnline)
	if err != nil {
		return nil, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	p := &player{
		ID:    ID,
		Alias: alias,
		State: PlayerOnline,
	}

	return p, nil
}

func (pr *PlayerRepository) Save(p *player) error {
	_, err := pr.NamedQuery(`UPDATE player SET state=:state WHERE alias=:alias`, p)

	return err
}

func (pr *PlayerRepository) FindByAlias(alias string) (*player, error) {
	p := &player{}
	if err := pr.Get(p, `SELECT * FROM player WHERE alias=?`, alias); err != nil {
		return nil, err
	}

	return p, nil
}
