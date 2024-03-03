package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pbaettig/moncron/internal/pkg/buildinfo"
	"github.com/pbaettig/moncron/internal/pkg/model"
	"github.com/pbaettig/moncron/internal/pkg/target"
	log "github.com/sirupsen/logrus"
)

var (
	logger *log.Entry
)

func gatherTargets(args *cmdlineArgs) []target.ResultTarget {
	targets := make([]target.ResultTarget, 0)
	if args.PushgatewayURL != "" {
		targets = append(targets, target.NewPrometheusPushgateway(args.PushgatewayURL))
	}

	if args.MonServerURL != "" {
		targets = append(targets, target.NewWebhook(args.MonServerURL))
	}

	if args.LogFile != "" {
		targets = append(targets, target.NewFile(args.LogFile, args.LogFileMaxSize))
	}

	if args.Stdout {
		targets = append(targets, target.Stdout{})
	}
	return targets

}

func parseArgs() *cmdlineArgs {
	args := new(cmdlineArgs)
	args.FromCmdline()

	if args.Version {
		// don't perform any further checks
		return args
	}

	if args.JobName == "" {
		flag.Usage()
		log.Fatalln("-name is required")
	}

	if len(args.ProcessCmdline) == 0 {
		flag.Usage()
		log.Fatalln("nothing to execute")
	}

	if args.Quiet {
		log.SetLevel(log.PanicLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return args
}

func main() {
	args := new(cmdlineArgs)
	// take values from the environment first
	args.FromEnv()
	// use values form cmdline second, giving them priority over
	// whatever was defined in the env vars
	args.FromCmdline()

	if err := args.Validate(); err != nil {
		log.Fatalln(err)
	}
	if args.Version {
		fmt.Printf(buildinfo.Version)
		os.Exit(0)
	}

	targets := gatherTargets(args)

	logger = log.WithFields(log.Fields{"name": args.JobName})

	job := &model.Job{Name: args.JobName}
	job.WithCommand(model.NewCommand(args.ProcessCmdline, args.Timeout))

	jobRun, err := job.Run()
	if err != nil {
		log.Errorf("Error running %s: %s", args.JobName, err)
	}

	for _, target := range targets {
		l := logger.WithField("target", target.Name())
		if err := target.Push(jobRun); err != nil {
			l.Warnf("could not push results to %s\n", err)
		} else {
			l.Infof("successfully pushed results")
		}
	}

	fmt.Println()
	fmt.Fprint(os.Stdout, jobRun.Result.Stdout)
	fmt.Fprint(os.Stderr, jobRun.Result.Stderr)
	os.Exit(0)
}
