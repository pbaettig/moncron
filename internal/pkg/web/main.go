package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pbaettig/moncron/internal/pkg/run"
)

func PushResults(name string, r run.CommandResult, url string) error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)

	return err
}
