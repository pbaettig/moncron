package main

import (
	"fmt"
	"sync"

	"github.com/pbaettig/moncron/internal/pkg/model"
	"github.com/pbaettig/moncron/internal/pkg/store/sqlite"
	"github.com/pbaettig/moncron/internal/pkg/testdata"
	log "github.com/sirupsen/logrus"
)

var (
	wg *sync.WaitGroup = new(sync.WaitGroup)
)

const (
	seed = 12454
)

func init() {
	testdata.Seed = seed
}

func populateDB(dbPath string, nRuns int, nHosts int) {
	log.SetLevel(log.InfoLevel)
	jobRunStore, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	for i, run := range testdata.GenerateJobRuns(nRuns, nHosts) {
		wg.Add(1)
		go func(i int, run model.JobRun) {
			defer wg.Done()

			if err := jobRunStore.Add(run); err != nil {
				fmt.Printf("%#v, %+v, %T\n", err, err, err)
				log.Error(err)
			}
			log.Infof("inserted jobRun #%d: %s,%s (%s)", i+1, run.ID, run.Name, run.Host.Name)
		}(i, run)
	}
	wg.Wait()
}

func main() {
	args := new(cmdlineArgs)
	args.Parse()
	populateDB(args.dbPath, args.numRuns, args.numHosts)
}
