package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pbaettig/moncron/internal/pkg/run"
)

type Webhook struct {
	ResultTarget
	url    string
	method string
}

type WebhookContent struct {
	JobName    string
	ExecutedAt string
	Result     run.CommandResult
}

func (w Webhook) Push(jobName string, r run.CommandResult) error {
	content := WebhookContent{
		JobName:    jobName,
		ExecutedAt: time.Now().UTC().Format("2006-01-02T15:04:05"),
		Result:     r,
	}

	body, err := json.Marshal(content)
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
		return fmt.Errorf("cannot push results to %s: %w", w.url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot push results to %s: got HTTP%s", w.url, resp.Status)
	}

	return err
}

func NewWebhook(url string) Webhook {
	return Webhook{
		url:    url,
		method: "POST",
	}
}
