package internal

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/safe-k/gonnect"
)

type storage struct {
	*pop.Connection
}

func (s *storage) GetActiveMatch(p gonnect.Player) (*gonnect.Match, error) {
	q := s.Q().
		LeftJoin("matches_players", "matches_players.match_id=matches.id").
		LeftJoin("players", "players.id=matches_players.player_id").
		Where("players.alias = ?", p.Alias).
		Where("matches.state <> ?", gonnect.MatchEnded)

	m := &gonnect.Match{}
	if err := q.First(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *storage) GetPlayersSearching() (*gonnect.Players, error) {
	pp := &gonnect.Players{}
	if err := s.Where("state = ?", gonnect.PlayerSearching).All(pp); err != nil {
		return nil, err
	}

	return pp, nil
}

func (s *storage) SavePlayer(p *gonnect.Player) error {
	return s.Save(p)
}

func (s *storage) SavePlayers(pp *gonnect.Players) error {
	return s.Save(pp)
}

func (s *storage) GetMatchById(id int) (*gonnect.Match, error) {
	m := &gonnect.Match{}
	if err := s.Find(m, id); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *storage) GetMatchesByState(state string) (*gonnect.Matches, error) {
	mm := &gonnect.Matches{}
	if err := s.Where("state = ?", state).All(mm); err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *storage) EndMatch(id int) error {
	return s.Update(&gonnect.Match{ID: id, State: gonnect.MatchEnded})
}

func (s *storage) SaveMatch(m *gonnect.Match) error {
	return s.Save(m)
}

func Storage() *storage {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalln("Could not connect to NewDB", err)
	}

	return &storage{db}
}
