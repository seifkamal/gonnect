package domain

import "time"

const (
	MatchReady    = "ready"
	MatchEnded    = "ended"
)

type Match struct {
	ID        int       `json:"id" db:"id"`
	State     string    `json:"state" db:"state"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Players   Players   `many_to_many:"matches_players" db:"-"`
}
