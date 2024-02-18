package model

import (
	"context"
	"os/exec"
	"time"
)

type Command struct {
	Executable string
	Args       []string
	ctx        context.Context
	cancel     context.CancelFunc
	cmd        *exec.Cmd
}

func NewCommand(args []string, timeout time.Duration) *Command {
	var (
		bgCtx context.Context = context.Background()
	)

	command := Command{
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
