package player

import (
	"github.com/safe-k/gonnect/internal"
)

const (
	Offline  = "offline"
	Online   = "online"
)

type player struct {
	ID    int64  `db:"id"`
	Alias string `db:"alias"`
	State string `db:"state"`
}

func New(alias string) *player {
	return &player{
		Alias: alias,
		State: "online",
	}
}

type repository internal.Repository

func Repository(storage internal.Storage) *repository {
	return &repository{Storage: storage}
}

func (r *repository) Insert(p *player) error {
	res, err := r.ExecuteTransaction(`INSERT INTO player (alias, state) VALUES(:alias, :state)`, p)
	if err != nil {
		return err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = ID

	return nil
}

func (r *repository) Update(p *player) error {
	_, err := r.ExecuteTransaction(`UPDATE player SET state=:state WHERE alias=:alias`, p)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByAlias(alias string) (*player, error) {
	p := &player{}
	if err := r.FindOne(p, `SELECT * FROM player WHERE alias=?`, alias); err != nil {
		return nil, err
	}

	return p, nil
}
