package model

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"syscall"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/sirupsen/logrus"
)

type Job struct {
	Name        string `validate:"required"`
	Description string
	Schedule    *cronexpr.Expression
	Command     *Command `validate:"required"`
}

func (j *Job) WithSchedule(cron string) *Job {
	j.Schedule = cronexpr.MustParse(cron)
	return j
}

func (j *Job) WithCommand(c *Command) *Job {
	j.Command = c
	return j
}

func (j *Job) WithDescription(desc string) *Job {
	j.Description = desc
	return j
}

func (j *Job) Run() (JobRun, error) {
	defer j.Command.cancel()

	run := JobRun{Job: *j, ID: uuid.NewString()}

	run.Environment = make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		run.Environment[kv[0]] = kv[1]
	}

	u, err := user.Current()
	if err != nil {
		return run, err
	}
	if u != nil {
		run.User = *u
	}

	h, err := GetHost()
	if err != nil {
		return run, err
	}
	run.Host = h

	j.Command.cmd = exec.CommandContext(j.Command.ctx, j.Command.Executable, j.Command.Args...)

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	j.Command.cmd.Stdout = stdout
	j.Command.cmd.Stderr = stderr

	err = j.Command.cmd.Start()
	if err != nil {
		return run, err
	}

	run.StartedAt = time.Now().UTC()
	logrus.Info("started")

	err = j.Command.cmd.Wait()

	run.FinishedAt = time.Now().UTC()
	run.Result.WallTime = run.FinishedAt.Sub(run.StartedAt)
	run.Result.ExitCode = j.Command.cmd.ProcessState.ExitCode()
	run.Result.Stdout = stdout.String()
	run.Result.Stderr = stderr.String()

	sys, ok := j.Command.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		run.Result.Killed = sys.Signal() == syscall.SIGKILL
		run.Result.ReceivedSignal = sys.Signal()
	}

	rusage, ok := j.Command.cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if ok {
		run.Result.SystemTime = time.Duration(rusage.Stime.Nano()) * time.Nanosecond
		run.Result.UserTime = time.Duration(rusage.Utime.Nano()) * time.Nanosecond
		run.Result.MaxRssBytes = rusage.Maxrss
	}

	return run, err
}

type Result struct {
	ExitCode       int
	SystemTime     time.Duration
	UserTime       time.Duration
	WallTime       time.Duration
	Killed         bool
	ReceivedSignal syscall.Signal
	MaxRssBytes    int64
	Stdout         string
	Stderr         string
}

type JobRun struct {
	Job
	ID          string `validate:"required,uuid"`
	Host        Host
	User        user.User
	StartedAt   time.Time
	FinishedAt  time.Time
	Environment map[string]string
	Result      Result
}
