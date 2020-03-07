package match

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	Creating = "creating"
	Ready    = "ready"
	Ended    = "ended"
)

type match struct {
	ID    int64  `db:"id"`
	State string `db:"state"`
}

type Repository struct {
	*sqlx.DB
}

func (r *Repository) FindByPlayerAlias(alias string) (*match, error) {
	stmt := "SELECT m.id, m.state FROM `match` AS m LEFT JOIN match_players AS mp ON mp.match_id=m.id LEFT JOIN player AS p ON p.id=mp.player_id WHERE m.state <> ? AND p.alias=?;"

	m := &match{}
	if err := r.Get(m, stmt, Ended, alias); err != nil {
		return nil, err
	}

	return m, nil
}

func (r *Repository) Find(state string) (*[]match, error) {
	mm := &[]match{}
	if err := r.Select(mm, "SELECT * FROM `match` WHERE state=?", state); err != nil {
		return nil, err
	}

	return mm, nil
}

func (r *Repository) Create() (*match, error) {
	res, err := r.Exec("INSERT INTO `match` (state) VALUES (?);", Creating)
	if err != nil {
		return nil, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	m := &match{
		ID:    ID,
		State: Creating,
	}

	return m, nil
}

func (r *Repository) Save(m *match) error {
	_, err := r.NamedExec("UPDATE `match` SET state=:state WHERE id=:id;", m)

	return err
}

func (r *Repository) LinkPlayers(m *match, ppID []int64) error {
	var vv []string
	for _, pID := range ppID {
		vv = append(vv, fmt.Sprintf("(%v, %v)", m.ID, pID))
	}

	_, err := r.Exec(`INSERT INTO match_players (match_id, player_id) VALUES ` + strings.Join(vv, ", "))
	return err
}
