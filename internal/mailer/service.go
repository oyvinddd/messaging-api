package mailer

import (
	"fmt"
	"bytes"
	"context"
	"net/http"
	"encoding/json"
)

const (
	postmarkAPI = "https://api.postmarkapp.com/email"
	postmarkMsgStream = "outbound"
)

type (
	Service interface {
		Send(ctx context.Context, message Message) error
	}

	postmarkService struct {
		apiKey string
	}
)

func NewPostmarkService(apiKey string) Service {
	return &postmarkService{apiKey: apiKey}
}

func (s *postmarkService) Send(ctx context.Context, message Message) error {
	// anonymous struct representing the http request body
	requestBody := struct{
		From          string `json:"From"`
		To            string `json:"To"`
		Subject       string `json:"Subject"`
		TextBody      string `json:"TextBody"`
		HTMLBody      string `json:"HtmlBody"`
		MessageStream string `json:"MessageStream"`
	}{
		From: message.Sender,
		To: message.Recipient,
		Subject: message.Subject,
		TextBody: message.Body,
		MessageStream: postmarkMsgStream,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, postmarkAPI, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.apiKey)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("postmark returned status %s", res.Status)
	}

	return nil
}

