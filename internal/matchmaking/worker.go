package matchmaking

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/safe-k/gonnect"
)

type (
	workerStorage interface {
		GetPlayersSearching() (*gonnect.Players, error)
		SaveMatch(match *gonnect.Match) error
		SavePlayers(players *gonnect.Players) error
	}

	worker struct {
		Storage workerStorage
	}
)

func (w *worker) Work(batch int) {
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

func (w *worker) createMatch(mpp gonnect.Players) {
	m := &gonnect.Match{
		State:   gonnect.MatchReady,
		Players: mpp,
	}

	if err := w.Storage.SaveMatch(m); err != nil {
		log.Fatalln("Could not create m", err)
	}

	if err := w.Storage.SavePlayers(mpp.Reserve()); err != nil {
		log.Fatalln("Could not reserve m players", err)
	}
}

func Worker(storage workerStorage) *worker {
	return &worker{
		Storage: storage,
	}
}
