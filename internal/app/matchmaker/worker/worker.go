package worker

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/safe-k/gonnect/internal/domain"
)

type storage interface {
	GetPlayersSearching() (*domain.Players, error)
	SaveMatch(match *domain.Match) error
	SavePlayers(players *domain.Players) error
}

type Worker struct {
	Storage storage
}

func (w *Worker) Work(bch int) {
	for {
		pp, err := w.Storage.GetPlayersSearching()
		if err != nil {
			if !strings.Contains(err.Error(), "no rows") {
				log.Fatalln("Could not find players")
			}
		}

		mc := len(*pp) / bch
		if mc == 0 {
			log.Println("Waiting for more players")
			<-time.After(2 * time.Second)
			continue
		}

		log.Println("Creating matches. Count:", mc)

		var (
			wg sync.WaitGroup
			s  = 0
			e  = bch
		)

		for c := 1; c <= mc; c++ {
			mpp := (*pp)[s:e]
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.createMatch(mpp)
			}()

			s += bch
			e += bch
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
		log.Fatalln("Could not create match", err)
	}

	if err := w.Storage.SavePlayers(mpp.Reserve()); err != nil {
		log.Fatalln("Could not reserve match players", err)
	}
}
