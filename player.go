package gonnect

import "time"

const (
	PlayerSearching   = "searching"
	PlayerUnavailable = "unavailable"
)

type (
	Player struct {
		ID        int       `json:"id" db:"id"`
		Alias     string    `json:"alias" db:"alias"`
		State     string    `json:"state" db:"state"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	Players []Player
)

func (pp Players) Reserve() *Players {
	for i := range pp {
		pp[i].State = PlayerUnavailable
	}

	return &pp
}
