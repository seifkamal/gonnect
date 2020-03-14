package internal

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/safe-k/gonnect/internal/domain"
)

func Storage() *storage {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalln("Could not connect to NewDB", err)
	}

	return &storage{db}
}

type storage struct {
	*pop.Connection
}

func (s *storage) GetActiveMatch(p domain.Player) (*domain.Match, error) {
	q := s.Q().
		LeftJoin("matches_players", "matches_players.match_id=matches.id").
		LeftJoin("players", "players.id=matches_players.player_id").
		Where("players.alias = ?", p.Alias).
		Where("matches.state <> ?", domain.MatchEnded)

	m := &domain.Match{}
	if err := q.First(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *storage) GetPlayersSearching() (*domain.Players, error) {
	pp := &domain.Players{}
	if err := s.Where("state = ?", domain.PlayerSearching).All(pp); err != nil {
		return nil, err
	}

	return pp, nil
}

func (s *storage) SavePlayer(p *domain.Player) error {
	return s.Save(p)
}

func (s *storage) SavePlayers(pp *domain.Players) error {
	return s.Save(pp)
}

func (s *storage) GetMatchById(id int) (*domain.Match, error) {
	m := &domain.Match{}
	if err := s.Find(m, id); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *storage) GetMatchesByState(state string) (*domain.Matches, error) {
	mm := &domain.Matches{}
	if err := s.Where("state = ?", state).All(mm); err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *storage) EndMatch(id int) error {
	return s.Update(&domain.Match{ID: id, State: domain.MatchEnded})
}

func (s *storage) SaveMatch(m *domain.Match) error {
	return s.Save(m)
}