package store

import (
	"errors"
	"time"

	"github.com/pbaettig/moncron/internal/pkg/model"
)

var (
	ErrAlreadyExists   = errors.New("job already exists")
	ErrNotFound        = errors.New("not found")
	ErrPageOutOfBounds = errors.New("requested page is out of bounds")
)

type JobRunStorer interface {
	Add(j model.JobRun) error
	GetHosts(page int, size int) (hosts []model.Host, total int, err error)
	GetHostNames(page int, size int) (hosts []string, total int, err error)
	GetJobNames(page int, size int) (jobs []string, total int, err error)
	GetByHost(hostName string, page int, size int) (runs []model.JobRun, total int, next int, err error)
	GetByJob(jobName string, page int, size int) (runs []model.JobRun, total int, next int, err error)
	GetByHostAndJob(jobName string, hostName string, page int, size int) (runs []model.JobRun, total int, next int, err error)
	GetByNameBefore(jobName string, t time.Time, page int, size int) (runs []model.JobRun, total int, next int, err error)
	GetStatsByHost() (stats map[string]*model.Counts, err error)
	Delete(id string) error
	Get(id string) (run model.JobRun, err error)
}
