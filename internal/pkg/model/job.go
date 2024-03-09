package model

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"syscall"

	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/sirupsen/logrus"
)

const (
	ProcessStartedNormally = "started-normally"
	ProceessNotStarted     = "start-denied"

	IOBlockSize = 512
)

type Counts struct {
	Successful int
	Failed     int
	Denied     int
	Total      int
}

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

func (j *Job) PrepareRun() *JobRun {
	r := new(JobRun)
	r.init(j)

	return r
}

type Result struct {
	ExitCode   int
	SystemTime time.Duration
	UserTime   time.Duration
	WallTime   time.Duration

	IOBlocksRead               uint64
	IOBytesRead                uint64
	IOBlocksWritten            uint64
	IOBytesWritten             uint64
	VoluntaryContextSwitches   uint64
	InvoluntaryContextSwitches uint64
	ReceivedSignal             syscall.Signal
	MaxRssBytes                uint64
	MemoryUtilization          float64
	Stdout                     string
	Stderr                     string
}

func (r Result) Killed() bool {
	return r.ReceivedSignal == syscall.SIGKILL
}

func (r Result) CPUUtilization() float64 {
	return ((float64(r.SystemTime) + float64(r.UserTime)) * 100.0) / float64(r.WallTime)
}

func (r Result) IdleTime() time.Duration {
	idle := r.WallTime - r.SystemTime - r.UserTime
	if idle < 0 {
		return time.Duration(0)
	}
	return idle
}

type JobRun struct {
	Job
	Status      string
	ID          string `validate:"required,uuid"`
	Host        Host
	User        user.User
	StartedAt   time.Time
	FinishedAt  time.Time
	Environment map[string]string
	Result      Result
}

func (r *JobRun) init(job *Job) error {
	if job == nil {
		return fmt.Errorf("job is nil")
	}
	r.Job = *job
	r.ID = uuid.NewString()

	return nil
}

func (r *JobRun) setHost() error {
	hi, err := host.Info()
	if err != nil {
		return err
	}
	r.Host.Name = hi.Hostname
	r.Host.OS.Name = hi.Platform
	r.Host.OS.KernelVersion = hi.KernelVersion
	r.Host.OS.Version = hi.PlatformVersion
	(&r.Host).SetCPU()
	(&r.Host).SetMemory()
	return nil
}

func (r *JobRun) setEnvironment() error {
	r.Environment = make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		r.Environment[kv[0]] = kv[1]
	}
	return nil
}

func (r *JobRun) setUser() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	if u != nil {
		r.User = *u
	}

	return nil
}

func (r *JobRun) NotRun() error {
	r.Status = ProceessNotStarted
	r.StartedAt = time.Now().UTC()
	r.FinishedAt = r.StartedAt
	if err := r.setEnvironment(); err != nil {
		return err
	}
	if err := r.setHost(); err != nil {
		return err
	}
	if err := r.setUser(); err != nil {
		return err
	}

	return nil
}

func (r *JobRun) Run() error {
	defer r.Job.Command.cancel()

	r.Command.cmd = exec.CommandContext(r.Job.Command.ctx, r.Job.Command.Executable, r.Job.Command.Args...)

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	r.Job.Command.cmd.Stdout = stdout
	r.Job.Command.cmd.Stderr = stderr

	if err := r.setEnvironment(); err != nil {
		return err
	}
	if err := r.setHost(); err != nil {
		return err
	}
	if err := r.setUser(); err != nil {
		return err
	}

	err := r.Job.Command.cmd.Start()
	if err != nil {
		return err
	}

	r.StartedAt = time.Now().UTC()

	r.Status = ProcessStartedNormally
	r.StartedAt = time.Now().UTC()
	logrus.Info("started")

	err = r.Job.Command.cmd.Wait()

	r.FinishedAt = time.Now().UTC()
	r.Result.WallTime = r.FinishedAt.Sub(r.StartedAt)
	r.Result.ExitCode = r.Job.Command.cmd.ProcessState.ExitCode()
	r.Result.Stdout = stdout.String()
	r.Result.Stderr = stderr.String()

	sys, ok := r.Job.Command.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		r.Result.ReceivedSignal = sys.Signal()
	}

	rusage, ok := r.Job.Command.cmd.ProcessState.SysUsage().(*syscall.Rusage)
	if ok {
		r.Result.MaxRssBytes = uint64(rusage.Maxrss) * 1024

		r.Result.IOBlocksRead = uint64(rusage.Inblock)
		r.Result.IOBytesRead = uint64(rusage.Inblock) * IOBlockSize

		r.Result.IOBlocksWritten = uint64(rusage.Oublock)
		r.Result.IOBytesWritten = uint64(rusage.Oublock) * IOBlockSize

		r.Result.InvoluntaryContextSwitches = uint64(rusage.Nivcsw)

		r.Result.VoluntaryContextSwitches = uint64(rusage.Nvcsw)

		r.Result.SystemTime = time.Duration(rusage.Stime.Nano()) * time.Nanosecond
		r.Result.UserTime = time.Duration(rusage.Utime.Nano()) * time.Nanosecond

		r.Result.MemoryUtilization = (float64(r.Result.MaxRssBytes) * 100) / float64(r.Host.MemoryBytes)
	}

	return err
}
