package main

import (
	"flag"
	"time"
)

type cmdlineArgs struct {
	Timeout        time.Duration
	JobName        string
	PushgatewayURL string
	WebhookURL     string
	Stdout         bool
	Quiet          bool
	Version        bool
	LogFile        string
	LogFileMaxSize int64
	ProcessCmdline []string
}

func (c *cmdlineArgs) Parse() error {
	flag.DurationVar(&c.Timeout, "timeout", time.Duration(0), "timeout value")
	flag.StringVar(&c.JobName, "name", "", "job name")
	flag.StringVar(&c.PushgatewayURL, "pushgw", "", "Prometheus pushgateway URL")
	flag.StringVar(&c.WebhookURL, "web", "", "Webhook URL")
	flag.BoolVar(&c.Stdout, "stdout", false, "echo results to stdout")
	flag.BoolVar(&c.Quiet, "quiet", false, "")
	flag.BoolVar(&c.Version, "version", false, "")
	flag.StringVar(&c.LogFile, "log", "", "Log File path")
	flag.Int64Var(&c.LogFileMaxSize, "log-size", 10*1024*1024, "Log File Maximum Size")
	flag.Parse()

	c.ProcessCmdline = flag.Args()
	return nil
}
