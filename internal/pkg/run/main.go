package run

import (
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type CommandResult struct {
	ExitCode       int
	SystemTimeNano int64
	UserTimeNano   int64
	WallTime       time.Duration
	Killed         bool
	ReceivedSignal syscall.Signal
	MaxRSS         int64
	Stdout         string
	Stderr         string
}

type Command struct {
	Name   string
	Args   []string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewCommand(args []string, timeout time.Duration) Command {
	var (
		bgCtx context.Context = context.Background()
	)

	command := Command{
		Name: args[0],
		Args: args[1:len(args)],
	}

	if timeout > 0 {
		command.ctx, command.cancel = context.WithTimeout(bgCtx, timeout)
	} else {
		command.ctx, command.cancel = context.WithCancel(bgCtx)
	}

	return command
}

func (c Command) Execute() (CommandResult, error) {
	var result CommandResult = CommandResult{}

	defer c.cancel()
	cmd := exec.CommandContext(c.ctx, c.Name, c.Args...)
	cmd.Stdout = new(strings.Builder)
	cmd.Stderr = new(strings.Builder)

	start := time.Now()
	err := cmd.Run()
	result.WallTime = time.Now().Sub(start)
	result.ExitCode = cmd.ProcessState.ExitCode()
	result.Stdout = cmd.Stdout.(*strings.Builder).String()
	result.Stderr = cmd.Stderr.(*strings.Builder).String()

	sys, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		result.Killed = sys.Signal() == syscall.SIGKILL
		result.ReceivedSignal = sys.Signal()
	}

	rusage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if ok {
		result.SystemTimeNano = rusage.Stime.Nano() / 1000 / 1000
		result.UserTimeNano = rusage.Utime.Nano() / 1000 / 1000
		result.MaxRSS = rusage.Maxrss
	}

	return result, err
}
