package receiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type slack struct {
	webhookUrl    string
	messageHeader string
}

type slackOption func(*slack)

func (s *slack) applyOptions(opts ...slackOption) {
	for _, o := range opts {
		o(s)
	}
}

func MessageHeader(msg string) slackOption {
	return func(s *slack) {
		s.messageHeader = msg
	}
}

func NewSlackReceiver(webhookUrl string, opts ...slackOption) (Notifier, error) {
	s := &slack{}
	s.messageHeader = "New GitHub issue"
	s.webhookUrl = webhookUrl
	s.applyOptions(opts...)
	return s, nil
}

func (s *slack) Notify(msg string) error {
	payload, err := json.Marshal(slackRequest{Text: msg})
	if err != nil {
		return fmt.Errorf("failed creating payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.webhookUrl, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to make http request: %v", err)
	}
	req.Header.Add("Content-Type", "applications/json")
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	buff := bytes.Buffer{}
	buff.ReadFrom(resp.Body)
	defer resp.Body.Close()
	if s := buff.String(); s != "ok" {
		return fmt.Errorf("failed connecting to Slack: %v", s)
	}
	return nil
}

type slackRequest struct {
	Text string `json:"text"`
}
