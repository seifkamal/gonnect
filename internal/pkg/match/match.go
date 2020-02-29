package match

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	Creating = "creating"
	Ready    = "ready"
	Ended    = "ended"
)

type match struct {
	ID    int    `db:"id"`
	State string `db:"state"`
}

type Repository struct {
	*sqlx.DB
}

func (r *Repository) FindByPlayerAlias(alias string) (*match, error) {
	// TODO: Find matches in state `creating`
	stmt := "SELECT m.id, m.state FROM `match` AS m LEFT JOIN match_players AS mp ON mp.match_id=m.id LEFT JOIN player AS p ON p.id=mp.player_id WHERE p.alias=?"

	m := &match{}
	if err := r.Get(m, stmt, alias); err != nil {
		return nil, err
	}

	return m, nil
}
