package target

import (
	"github.com/pbaettig/moncron/internal/pkg/run"
)

type ResultTarget interface {
	Push(name string, r run.CommandResult) error
}
