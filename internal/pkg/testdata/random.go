package testdata

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pbaettig/moncron/internal/pkg/model"
)

func RandomHost() model.Host {
	h := model.Host{}
	h.Name = fmt.Sprintf("srv-%d.%s", gofakeit.IntRange(1000, 9999), gofakeit.DomainName())
	h.OS = model.OperatingSystem{
		Name:          "debian",
		Version:       "sid",
		KernelVersion: "2.7.23",
	}
	return h
}

func RandomCommand(executableName string) *model.Command {
	c := new(model.Command)
	c.Executable = strings.ToLower(executableName)
	c.Args = make([]string, 0)
	return c
}

func RandomJob() model.Job {
	j := model.Job{Command: new(model.Command)}
	j.Name = fmt.Sprintf("./%s-%d", gofakeit.AppName(), gofakeit.Int())
	j.Description = gofakeit.HackerPhrase()
	j.Command = RandomCommand(j.Name)

	return j
}

func RandomResult() model.Result {
	r := model.Result{}
	r.ExitCode = gofakeit.IntRange(0, 10)
	r.SystemTime = time.Duration(gofakeit.IntRange(500, 5000)) * time.Second
	r.UserTime = time.Duration(gofakeit.IntRange(500, 5000)) * time.Second
	r.WallTime = r.SystemTime + r.UserTime
	r.Killed = gofakeit.IntRange(1, 10) == 10
	if r.Killed {
		r.ReceivedSignal = syscall.SIGKILL
		r.ExitCode = -1
	}
	r.MaxRssBytes = int64(gofakeit.IntRange(10000, 400000000))
	if r.ExitCode > 0 {
		for i := 0; i < gofakeit.IntRange(1, 20); i++ {
			r.Stderr += gofakeit.Sentence(gofakeit.IntRange(4, 10)) + "\n"
		}
		r.Stderr += fmt.Sprintf("\n\n%s", gofakeit.Error())

	} else {
		r.Stdout = gofakeit.HackerNoun()
	}

	return r
}

func randomPastDate() time.Time {
	now := time.Now()
	d := -time.Duration(gofakeit.IntRange(24*3600, 3*12*30*24*3600)) * time.Second

	return now.Add(d)

}

func RandomJobRun() model.JobRun {
	j := RandomJob()
	r := j.PrepareRun()

	if gofakeit.IntN(20) == 19 { // ~5%
		r.NotRun()
	}

	r.Result = RandomResult()
	r.Command = RandomCommand(r.Job.Command.Executable)
	r.ID = gofakeit.UUID()
	r.StartedAt = randomPastDate()
	r.FinishedAt = r.StartedAt.Add(r.Result.WallTime)
	r.Environment = make(map[string]string)
	r.Host = RandomHost()
	r.User = user.User{
		Uid:      strconv.Itoa(gofakeit.IntRange(100, 999)),
		Gid:      strconv.Itoa(gofakeit.IntRange(100, 999)),
		Username: strings.ToLower(gofakeit.FirstName()),
	}

	return *r
}
