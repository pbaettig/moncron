package target

import (
	"github.com/pbaettig/moncron/internal/pkg/run"
)

type ResultTarget interface {
	Push(r *run.Command) error
	Name() string
}
