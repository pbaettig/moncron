package main

import (
	"flag"
	"fmt"
)

type cmdlineArgs struct {
	dbPath   string
	numRuns  int
	numHosts int
}

func (c *cmdlineArgs) Parse() error {
	flag.StringVar(&c.dbPath, "db", "", "Path to sqlite DB file")
	flag.IntVar(&c.numRuns, "runs", 100, "number of job runs to generate")
	flag.IntVar(&c.numHosts, "hosts", 15, "number of hosts to distribute runs across")
	flag.Parse()

	if c.dbPath == "" {
		return fmt.Errorf("-db is mandatory")
	}

	return nil
}
