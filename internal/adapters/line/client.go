package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrMessengerNotConfigured = errors.New("line messenger is not configured")

type Messenger interface {
	ReplyText(replyToken, text string) error
	PushText(to, text string) error
}

type Client struct {
	channelAccessToken string
	httpClient         *http.Client
}

func NewClient(channelAccessToken string) *Client {
	return &Client{
		channelAccessToken: channelAccessToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) ReplyText(replyToken, text string) error {
	if c.channelAccessToken == "" || replyToken == "" || text == "" {
		return nil
	}

	payload := map[string]any{
		"replyToken": replyToken,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": text,
			},
		},
	}
	return c.sendMessage("https://api.line.me/v2/bot/message/reply", payload)
}

func (c *Client) PushText(to, text string) error {
	if c.channelAccessToken == "" {
		return ErrMessengerNotConfigured
	}
	if to == "" || text == "" {
		return nil
	}

	payload := map[string]any{
		"to": to,
		"messages": []map[string]string{
			{
				"type": "text",
				"text": text,
			},
		},
	}
	return c.sendMessage("https://api.line.me/v2/bot/message/push", payload)
}

func (c *Client) sendMessage(url string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.channelAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("line message request failed with status %d", resp.StatusCode)
	}
	return nil
}
