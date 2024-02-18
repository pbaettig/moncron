package main

import (
	"flag"
	"time"
)

type cmdlineArgs struct {
	listenAddress   string
	listenPort      int
	dbPath          string
	shutdownTimeout time.Duration
}

func (c *cmdlineArgs) Parse() error {
	flag.StringVar(&c.listenAddress, "listen", "0.0.0.0", "Listen address")
	flag.IntVar(&c.listenPort, "port", 8088, "Listen port")
	flag.StringVar(&c.dbPath, "db", "test.db", "Path to sqlite DB file")
	flag.DurationVar(&c.shutdownTimeout, "gtimeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	return nil
}
