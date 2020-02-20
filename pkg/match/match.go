package match

import "github.com/safe-k/gonnect/internal"

type match struct {
	ID    int    `db:"id"`
	State string `db:"state"`
}

type repository internal.Repository

func Repository(storage internal.Storage) *repository {
	return &repository{Storage: storage}
}

func (r *repository) FindByPlayerAlias(alias string) (*match, error) {
	// TODO: Find matches in state `creating`
	stmt := "SELECT m.id, m.state FROM `match` AS m LEFT JOIN match_players AS mp ON mp.match_id=m.id LEFT JOIN player AS p ON p.id=mp.player_id WHERE p.alias=?"

	m := &match{}
	if err := r.FindOne(m, stmt, alias); err != nil {
		return nil, err
	}

	return m, nil
}
