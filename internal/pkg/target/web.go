package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pbaettig/moncron/internal/pkg/model"
)

type Webhook struct {
	ResultTarget
	url    string
	method string
}

// type WebhookContent struct {
// 	Command    run.Command
// 	ExecutedAt string
// 	Result     run.CommandResult
// }

func (w Webhook) Name() string {
	return fmt.Sprintf("%s-webhook", w.method)
}

func (w Webhook) Push(r model.JobRun) error {

	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(w.method, w.url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "moncron")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP Error: %s", resp.Status)
	}

	return err
}

func NewWebhook(url string) Webhook {
	return Webhook{
		url:    url,
		method: "POST",
	}
}
