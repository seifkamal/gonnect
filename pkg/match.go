package pkg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type match struct {
	ID    int    `db:"id"`
	State string `db:"state"`
}

type MatchRepository struct {
	*sqlx.DB
}

func (mr *MatchRepository) FindByPlayerAlias(alias string) (*match, error) {
	// TODO: Find matches in state `creating`
	stmt := "SELECT m.id, m.state FROM `match` AS m LEFT JOIN match_players AS mp ON mp.match_id=m.id LEFT JOIN player AS p ON p.id=mp.player_id WHERE p.alias=?"

	m := &match{}
	if err := mr.Get(m, stmt, alias); err != nil {
		return nil, err
	}

	return m, nil
}
