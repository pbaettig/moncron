package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pbaettig/moncron/internal/pkg/run"
	"github.com/pbaettig/moncron/internal/pkg/target"
)

type cmdlineArgs struct {
	Timeout        time.Duration
	JobName        string
	PushgatewayURL string
	WebhookURL     string
	Verbose        bool
	ProcessCmdline []string
}

func (c *cmdlineArgs) Parse() error {
	flag.DurationVar(&c.Timeout, "timeout", time.Duration(0), "timeout value")
	flag.StringVar(&c.JobName, "name", "", "job name")
	flag.StringVar(&c.PushgatewayURL, "pushgw", "http://localhost:9091", "Prometheus pushgateway URL")
	flag.StringVar(&c.WebhookURL, "web", "http://localhost:8080", "Webhook URL")
	flag.BoolVar(&c.Verbose, "verbose", false, "")
	flag.Parse()
	if c.JobName == "" {
		return errors.New("-name cannot be empty")
	}

	c.ProcessCmdline = flag.Args()
	return nil
}

func main() {
	exitCode := 0

	args := new(cmdlineArgs)
	if err := args.Parse(); err != nil {
		flag.Usage()
		os.Exit(1)
	}

	if len(args.ProcessCmdline) == 0 {
		flag.Usage()
		fmt.Println("nothing to execute")
		os.Exit(1)
	}

	targets := make([]target.ResultTarget, 0)

	if args.PushgatewayURL != "" {
		targets = append(targets, target.NewPrometheusPushgateway(args.PushgatewayURL))
	}

	if args.WebhookURL != "" {
		targets = append(targets, target.NewWebhook(args.WebhookURL))
	}

	fmt.Printf("running %s...\n", args.JobName)
	cmd := run.NewCommand(args.ProcessCmdline, args.JobName, args.Timeout)
	result, err := cmd.Execute()

	fmt.Fprint(os.Stdout, result.Stdout)
	fmt.Fprint(os.Stderr, result.Stderr)
	exitCode = result.ExitCode

	for _, target := range targets {
		if err := target.Push(args.JobName, result); err != nil {
			fmt.Println(err)
		}
	}

	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			fmt.Printf("command failed with Exit code: %d\n", err.(*exec.ExitError).ExitCode())
			fmt.Printf("Error was: %s", err.Error())
		default:
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("deadline exceeded, increase timeout")
			}

		}

	}

	os.Exit(exitCode)
}
