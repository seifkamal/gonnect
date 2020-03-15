package matchmaking

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/safe-k/gonnect"
)

type mockWorkerStorage struct {
	players      gonnect.Players
	savedMatches chan gonnect.Match
	savedPlayers chan gonnect.Players
}

func (m *mockWorkerStorage) GetPlayersSearching() (*gonnect.Players, error) {
	pp := gonnect.Players{}
	for _, p := range m.players {
		if p.State == gonnect.PlayerSearching {
			pp = append(pp, p)
		}
	}

	if len(pp) == 0 {
		return nil, errors.New("no players found")
	}

	return &pp, nil
}

func (m *mockWorkerStorage) SaveMatch(match *gonnect.Match) error {
	m.savedMatches <- *match
	return nil
}

func (m *mockWorkerStorage) SavePlayers(players *gonnect.Players) error {
	m.savedPlayers <- *players
	return nil
}

func TestWorker(t *testing.T) {
	tests := []struct {
		batch, playerCount, matchCount int
	}{
		{2, 10, 5},
		{3, 10, 3},
		{4, 10, 2},
		{5, 10, 2},
		{10, 10, 1},
		{10, 20, 2},
		{5, 37, 7},
	}

	now := time.Now()
	initialState := gonnect.PlayerSearching
	finalState := gonnect.PlayerUnavailable
	for _, test := range tests {
		mockStorage := &mockWorkerStorage{}
		for i := 1; len(mockStorage.players) < test.playerCount; i++ {
			mockStorage.players = append(mockStorage.players, gonnect.Player{
				ID:        i,
				Alias:     "toast" + strconv.Itoa(i),
				State:     initialState,
				CreatedAt: now,
				UpdatedAt: now,
			})

			mockStorage.savedMatches = make(chan gonnect.Match, test.matchCount)
			mockStorage.savedPlayers = make(chan gonnect.Players, test.playerCount)
		}

		if err := Worker(mockStorage).Work(test.batch); err != nil {
			t.Error(err)
		}

		close(mockStorage.savedMatches)
		close(mockStorage.savedPlayers)

		matchCount := len(mockStorage.savedMatches)
		if matchCount != test.matchCount {
			t.Errorf("Expected %v matches but %v were created", test.matchCount, matchCount)
		}

		playerCount := len(mockStorage.players)
		if playerCount != test.playerCount {
			t.Errorf("Expected %v players but %v were created", test.playerCount, playerCount)
		}

		for savedPlayers := range mockStorage.savedPlayers {
			for _, player := range savedPlayers {
				if player.State != finalState {
					t.Errorf("Expected player state to be %s but found %s instead", finalState, player.State)
				}
			}
		}
	}
}
