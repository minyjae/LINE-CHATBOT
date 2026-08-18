package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/types"
)

func (h *lineController) Webhook(c *fiber.Ctx) error {
	body := c.BodyRaw()
	signature := c.Get("X-Line-Signature")

	if h.channelSecret != "" && !verifySignature(body, signature, h.channelSecret) {
		return h.Response.Unauthorized(c, "Invalid LINE signature", "INVALID_LINE_SIGNATURE")
	}

	var payload types.LineWebhookRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		return h.Response.BadRequest(c, "Invalid LINE webhook payload", "INVALID_LINE_PAYLOAD")
	}

	processed := 0
	for _, event := range payload.Events {
		if event.Type != "message" || event.Message.Type != "text" || event.Source.UserID == "" {
			continue
		}

		result, err := h.lineWebhookService.HandleTextMessage(types.LineTextMessageInput{
			LineUserID: event.Source.UserID,
			ReplyToken: event.ReplyToken,
			Text:       event.Message.Text,
			Now:        time.Now(),
		})
		if err != nil {
			return h.Response.InternalServerError(c, "Failed to handle LINE message", err.Error(), "LINE_WEBHOOK_FAILED")
		}

		if err := h.lineMessenger.ReplyText(event.ReplyToken, result.ReplyText); err != nil {
			return h.Response.InternalServerError(c, "Failed to reply LINE message", err.Error(), "LINE_REPLY_FAILED")
		}
		processed++
	}

	return h.Response.Item(c, "LINE webhook processed", fiber.Map{"processed": processed})
}

func verifySignature(body []byte, signature, channelSecret string) bool {
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
