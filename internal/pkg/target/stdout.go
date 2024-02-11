package target

import (
	"encoding/json"
	"fmt"

	"github.com/pbaettig/moncron/internal/pkg/run"
)

type Stdout struct{}

func (p Stdout) Name() string {
	return "stdout"
}

func (p Stdout) Push(r *run.Command) error {
	if r == nil {
		return fmt.Errorf("nothing to push")
	}

	buf, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(buf))

	return nil
}
