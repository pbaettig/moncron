package target

import (
	"github.com/pbaettig/moncron/internal/pkg/model"
)

type ResultTarget interface {
	Push(r model.JobRun) error
	Name() string
}
