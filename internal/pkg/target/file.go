package target

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pbaettig/moncron/internal/pkg/model"
)

type File struct {
	ResultTarget
	path         string
	maxSizeBytes int64
}

func (f File) Name() string {
	return fmt.Sprintf("file-%s", f.path)
}

func NewFile(path string, maxSize int64) File {
	return File{
		path:         path,
		maxSizeBytes: maxSize,
	}
}

func (f File) Push(r model.JobRun) error {

	fd, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		return err
	}

	line, err := json.Marshal(r)

	if stat.Size()+int64(len(line)) > f.maxSizeBytes {
		if err := fd.Truncate(0); err != nil {
			return err
		}
	}

	buf, err := json.Marshal(r)
	if err != nil {
		return err
	}

	buf = append(buf, '\n')
	_, err = fd.Write(buf)

	return err
}
