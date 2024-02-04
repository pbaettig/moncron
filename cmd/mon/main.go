package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"time"

	"github.com/pbaettig/moncron/internal/pkg/job"
	_ "github.com/pbaettig/moncron/internal/pkg/job"
	"github.com/pbaettig/moncron/internal/pkg/run"
)

type cmdlineArgs struct {
	Timeout        time.Duration
	ProcessCmdline []string
}

func (c *cmdlineArgs) Parse() {
	flag.DurationVar(&c.Timeout, "timeout", time.Duration(0), "timeout value")
	flag.Parse()
	c.ProcessCmdline = flag.Args()
}

func main() {
	fmt.Println("Hello")
	args := new(cmdlineArgs)
	args.Parse()

	fmt.Println(args.Timeout)
	fmt.Println(args.ProcessCmdline)

	cmd := run.NewCommand(args.ProcessCmdline, args.Timeout)
	result, err := cmd.Execute()

	fmt.Printf("%+v\n", result)

	if err != nil {

		fmt.Printf("%v\n", err)
		fmt.Printf("%T\n", err)

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

	fmt.Println("files:")
	job.FindAll()
}
