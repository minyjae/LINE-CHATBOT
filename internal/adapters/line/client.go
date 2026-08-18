package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Messenger interface {
	ReplyText(replyToken, text string) error
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
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/v2/bot/message/reply", bytes.NewReader(body))
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
		return fmt.Errorf("line reply failed with status %d", resp.StatusCode)
	}
	return nil
}
