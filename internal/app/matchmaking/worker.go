package matchmaking

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/safe-k/gonnect/internal/domain"
)

type workerStorage interface {
	GetPlayersSearching() (*domain.Players, error)
	SaveMatch(match *domain.Match) error
	SavePlayers(players *domain.Players) error
}

type Worker struct {
	Storage workerStorage
}

func (w *Worker) Work(batch int) {
	for {
		pp, err := w.Storage.GetPlayersSearching()
		if err != nil {
			if !strings.Contains(err.Error(), "no rows") {
				log.Fatalln("Could not find players")
			}
		}

		mCount := len(*pp) / batch
		if mCount == 0 {
			log.Println("Waiting for more players")
			<-time.After(2 * time.Second)
			continue
		}

		log.Println("Creating matches. Count:", mCount)

		var (
			wg sync.WaitGroup
			s  = 0
			e  = batch
		)

		for c := 1; c <= mCount; c++ {
			mpp := (*pp)[s:e]
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.createMatch(mpp)
			}()

			s += batch
			e += batch
		}

		wg.Wait()
	}
}

func (w *Worker) createMatch(mpp domain.Players) {
	m := &domain.Match{
		State:   domain.MatchReady,
		Players: mpp,
	}

	if err := w.Storage.SaveMatch(m); err != nil {
		log.Fatalln("Could not create m", err)
	}

	if err := w.Storage.SavePlayers(mpp.Reserve()); err != nil {
		log.Fatalln("Could not reserve m players", err)
	}
}
