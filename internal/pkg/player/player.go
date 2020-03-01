package player

import (
	"strconv"

	"github.com/jmoiron/sqlx"
)

const (
	Away      = "away"
	Searching = "searching"
	Reserved  = "reserved"
)

type player struct {
	ID    int64  `db:"id"`
	Alias string `db:"alias"`
	State string `db:"state"`
}

type Repository struct {
	*sqlx.DB
}

func (r *Repository) New(alias string, state string) (*player, error) {
	res, err := r.Exec(`INSERT INTO player (alias, state) VALUES(?, ?);`, alias, state)
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
		State: state,
	}

	return p, nil
}

func (r *Repository) Save(p *player) error {
	_, err := r.NamedExec(`UPDATE player SET state=:state WHERE alias=:alias;`, p)

	return err
}

func (r *Repository) Reserve(pp []player) error {
	var ppID []string
	for _, p := range pp {
		ppID = append(ppID, strconv.FormatInt(p.ID, 10))
	}

	query, args, err := sqlx.In(`UPDATE player SET state=? WHERE id IN (?);`, Reserved, ppID)
	if err != nil {
		return err
	}

	_, err = r.Query(r.Rebind(query), args...)
	return err
}

func (r *Repository) FindByAlias(alias string) (*player, error) {
	p := &player{}
	if err := r.Get(p, `SELECT * FROM player WHERE alias=?;`, alias); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repository) FindAllSearching() ([]player, error) {
	pp := &[]player{}
	if err := r.Select(pp, `SELECT * FROM player WHERE state=?;`, Searching); err != nil {
		return nil, err
	}

	return *pp, nil
}
