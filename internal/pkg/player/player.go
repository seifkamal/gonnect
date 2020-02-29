package player

import (
	"github.com/safe-k/gonnect/internal/pkg/database"
)

const (
	Offline   = "offline"
	Searching = "searching"
	Reserved  = "reserved"
	Engaged   = "engaged"
)

type player struct {
	ID    database.ID `db:"id"`
	Alias string      `db:"alias"`
	State string      `db:"state"`
}

type Repository database.Repository

func (r *Repository) New(alias string, state string) (*player, error) {
	res, err := r.Exec(`INSERT INTO player (alias, state) VALUES(?, ?)`, alias, state)
	if err != nil {
		return nil, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	p := &player{
		ID:    database.ID(ID),
		Alias: alias,
		State: state,
	}

	return p, nil
}

func (r *Repository) Save(p *player) error {
	_, err := r.NamedExec(`UPDATE player SET state=:state WHERE alias=:alias`, p)

	return err
}

func (r *Repository) FindByAlias(alias string) (*player, error) {
	p := &player{}
	if err := r.Get(p, `SELECT * FROM player WHERE alias=?`, alias); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) FindAllSearching() (*[]player, error) {
	pp := &[]player{}
	if err := r.Select(pp, `SELECT * FROM player WHERE state=?`, Searching); err != nil {
		return nil, err
	}

	return pp, nil
}
