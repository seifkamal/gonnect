package matcher

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/safe-k/gonnect/internal/pkg/database"
	"github.com/safe-k/gonnect/internal/pkg/match"
	"github.com/safe-k/gonnect/internal/pkg/player"
)

func Work(bch int) {
	db := database.New()
	defer db.Close()

	pr := player.Repository{DB: db}
	mr := match.Repository{DB: db}

	for {
		pp, err := pr.FindAllSearching()
		if err != nil && err != sql.ErrNoRows {
			log.Fatalln("Could not find players")
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
				m, err := mr.Create()
				if err != nil {
					log.Fatalln("Could not create match", err)
				}

				log.Println("Reserving players")
				if err := pr.Reserve(mpp); err != nil {
					log.Fatalln("Could not reserve players", err)
				}

				var mppID []int64
				for _, p := range mpp {
					mppID = append(mppID, p.ID)
				}

				log.Println("Linking players to match")
				if err := mr.LinkPlayers(m, mppID); err != nil {
					log.Fatalln("Could not link players to match", err)
				}

				log.Println("Updating match state")
				m.State = match.Ready
				if err := mr.Save(m); err != nil {
					log.Fatalln("Could not update match state", err)
				}
			}()

			s += bch
			e += bch
		}

		wg.Wait()
	}
}
