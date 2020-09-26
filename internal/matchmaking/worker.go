package matchmaking

import (
	"log"
	"sync"
	"time"

	"github.com/seifkamal/gonnect"
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

func (w *worker) WorkIndefinitely(batch, retryInterval int) {
	interval := time.Duration(retryInterval)
	for {
		if err := w.Work(batch); err != nil {
			switch err.(type) {
			case gonnect.NoResultsFound, gonnect.InsufficientPlayers:
				log.Println("Waiting for more players")
				<-time.After(interval * time.Second)
				continue
			default:
				log.Fatal(err)
			}
		}
	}
}

func (w *worker) Work(batch int) error {
	pp, err := w.Storage.GetPlayersSearching()
	if err != nil {
		return err
	}

	pCount := len(*pp)
	mCount := pCount / batch
	if mCount == 0 {
		return gonnect.InsufficientPlayers{
			Need: batch,
			Have: pCount,
		}
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
	return nil
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
