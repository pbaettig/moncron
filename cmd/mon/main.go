package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pbaettig/moncron/internal/pkg/buildinfo"
	"github.com/pbaettig/moncron/internal/pkg/run"
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

	if args.WebhookURL != "" {
		targets = append(targets, target.NewWebhook(args.WebhookURL))
	}

	if args.LogFile != "" {
		targets = append(targets, target.NewFile(args.LogFile, args.LogFileMaxSize))
	}

	return targets

}

func runCommand(cmd *run.Command) {

	err := cmd.Start()
	if err != nil {
		logger.Errorf("unable to start: %s", err)
	}
	logger.Infoln("started")

	err = cmd.Wait()

	var exitStatus string
	if cmd.Result.Killed {
		exitStatus = "killed"
	} else {
		exitStatus = strconv.Itoa(cmd.Result.ExitCode)
	}
	l := logger.WithField("exit", exitStatus)

	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			l.Warningf("command failed: %+v", e)
			e.Sys()
		default:
			l.Error(err)
			os.Exit(1)
		}
	} else {
		l.Infoln("command finished successfully")
	}
}

func parseArgs() *cmdlineArgs {
	args := new(cmdlineArgs)
	args.Parse()

	if args.Version {
		// don't perform any further checks
		return args
	}

	if len(args.ProcessCmdline) == 0 {
		flag.Usage()
		fmt.Println("nothing to execute")
		os.Exit(1)
	}

	if args.Quiet {
		log.SetLevel(log.PanicLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return args
}

func main() {
	args := parseArgs()
	if args.Version {
		fmt.Printf(buildinfo.Version)
		os.Exit(0)
	}

	logger = log.WithFields(log.Fields{"name": args.JobName})

	targets := gatherTargets(args)

	cmd := run.NewCommand(args.ProcessCmdline, args.JobName, args.Timeout)
	runCommand(cmd)

	for _, target := range targets {
		l := logger.WithField("target", target.Name())
		if err := target.Push(cmd); err != nil {
			l.Warnf("could not push results to %s\n", err)
		} else {
			l.Infof("successfully pushed results")
		}
	}

	fmt.Println()
	fmt.Fprint(os.Stdout, cmd.Result.Stdout)
	fmt.Fprint(os.Stderr, cmd.Result.Stderr)
	os.Exit(0)
}
