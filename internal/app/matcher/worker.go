package matcher

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/safe-k/gonnect/internal/domain"
)

func Work(bch int) {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalln("Could not connect to DB", err)
	}
	defer db.Close()

	for {
		var pp domain.Players
		if err := db.Where("state = ?", domain.PlayerSearching).All(&pp); err != nil {
			if !strings.Contains(err.Error(), "no rows") {
				log.Fatalln("Could not find players")
			}
		}

		mc := len(pp) / bch
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
			mpp := pp[s:e]
			wg.Add(1)
			go func() {
				defer wg.Done()
				m := &domain.Match{
					State:   domain.MatchReady,
					Players: mpp,
				}

				if err := db.Save(m); err != nil {
					log.Fatalln("Could not create match", err)
				}

				if err := db.Update(mpp.Reserve()); err != nil {
					log.Fatalln("Could not reserve match players", err)
				}
			}()

			s += bch
			e += bch
		}

		wg.Wait()
		break
	}
}
