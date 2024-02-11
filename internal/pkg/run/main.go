package run

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
	"time"
)

type Command struct {
	Name       string
	Executable string
	Args       []string
	Result     CommandResult
	ctx        context.Context
	cancel     context.CancelFunc
	cmd        *exec.Cmd
}

type CommandResult struct {
	User             *user.User
	Environment      map[string]string
	WorkingDirectory string
	StartedAt        time.Time
	FinishedAt       time.Time
	ExitCode         int
	SystemTime       syscall.Timeval
	UserTime         syscall.Timeval
	WallTime         time.Duration
	Killed           bool
	ReceivedSignal   syscall.Signal
	MaxRssBytes      int64
	Stdout           string
	Stderr           string
}

func (c CommandResult) MarshalJSON() ([]byte, error) {
	type Alias CommandResult

	signal := ""
	if c.ReceivedSignal > -1 {
		signal = c.ReceivedSignal.String()
	}

	return json.Marshal(struct {
		Alias
		WallTime       string
		UserTime       string
		SystemTime     string
		ReceivedSignal string
	}{
		WallTime:       c.WallTime.String(),
		UserTime:       (time.Nanosecond * time.Duration(c.UserTime.Nano())).String(),
		SystemTime:     (time.Nanosecond * time.Duration(c.SystemTime.Nano())).String(),
		ReceivedSignal: signal,
		Alias:          (Alias)(c),
	})

}

func NewCommand(args []string, name string, timeout time.Duration) *Command {
	var (
		bgCtx context.Context = context.Background()
	)

	command := Command{
		Name:       name,
		Executable: args[0],
		Args:       args[1:],
	}

	if timeout > 0 {
		command.ctx, command.cancel = context.WithTimeout(bgCtx, timeout)
	} else {
		command.ctx, command.cancel = context.WithCancel(bgCtx)
	}

	return &command
}

func (c *Command) Start() error {
	c.cmd = exec.CommandContext(c.ctx, c.Executable, c.Args...)
	c.cmd.Stdout = new(strings.Builder)
	c.cmd.Stderr = new(strings.Builder)

	u, _ := user.Current()

	wd, _ := os.Getwd()
	if c.cmd.Dir != "" {
		wd = c.cmd.Dir
	}

	c.Result = CommandResult{
		StartedAt:        time.Now().UTC(),
		User:             u,
		Environment:      make(map[string]string),
		WorkingDirectory: wd,
	}
	for _, e := range c.cmd.Environ() {
		kv := strings.SplitN(e, "=", 2)
		c.Result.Environment[kv[0]] = kv[1]
	}

	return c.cmd.Start()
}

func (c *Command) Wait() error {
	defer c.cancel()

	err := c.cmd.Wait()

	c.Result.FinishedAt = time.Now().UTC()
	c.Result.WallTime = c.Result.FinishedAt.Sub(c.Result.StartedAt)
	c.Result.ExitCode = c.cmd.ProcessState.ExitCode()
	c.Result.Stdout = c.cmd.Stdout.(*strings.Builder).String()
	c.Result.Stderr = c.cmd.Stderr.(*strings.Builder).String()

	sys, ok := c.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		c.Result.Killed = sys.Signal() == syscall.SIGKILL
		c.Result.ReceivedSignal = sys.Signal()

	}

	rusage, ok := c.cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if ok {
		c.Result.SystemTime = rusage.Stime
		c.Result.UserTime = rusage.Utime
		c.Result.MaxRssBytes = rusage.Maxrss
	}

	return err
}
