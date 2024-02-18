package target

import (
	"encoding/json"
	"fmt"

	"github.com/pbaettig/moncron/internal/pkg/model"
)

type Stdout struct{}

func (p Stdout) Name() string {
	return "stdout"
}

func (p Stdout) Push(r model.JobRun) error {
	buf, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(buf))

	return nil
}
