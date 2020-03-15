package matchmaker

import (
	"github.com/safe-k/gonnect/internal"
	"github.com/safe-k/gonnect/internal/app/matchmaker/worker"
)

func Match(bch int) {
	storage := internal.Storage()
	defer storage.Close()

	w := &worker.Worker{Storage: storage}
	w.Work(bch)
}
