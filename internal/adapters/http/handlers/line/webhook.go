package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"minyjae/go-starter/types"

	"github.com/gofiber/fiber/v2"
)

// Webhook รับ request จาก LINE Messaging API แล้วตอบกลับข้อความ
// input: Fiber context ที่มี raw body, X-Line-Signature header และ LINE webhook payload
// output: HTTP response จำนวน event ที่ประมวลผล หรือ error response
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

		// ถ้า event เป็น message หรือ text ก็จะ handleMassage ตาม text นั้นๆที่ส่งมา
		result, err := h.lineWebhookService.HandleTextMessage(types.LineTextMessageInput{
			LineUserID: event.Source.UserID,
			ReplyToken: event.ReplyToken,
			Text:       event.Message.Text,
			Now:        time.Now(),
		})
		if err != nil {
			return h.Response.InternalServerError(c, "Failed to handle LINE message", err.Error(), "LINE_WEBHOOK_FAILED")
		}

		// ตอบกลับข้อความไปยัง line
		if err := h.lineMessenger.ReplyText(event.ReplyToken, result.ReplyText); err != nil {
			return h.Response.InternalServerError(c, "Failed to reply LINE message", err.Error(), "LINE_REPLY_FAILED")
		}
		processed++
	}

	return h.Response.Item(c, "LINE webhook processed", fiber.Map{"processed": processed})
}

// verifySignature ตรวจความถูกต้องของ LINE signature ด้วย HMAC-SHA256
// input: body raw request, signature จาก header, channelSecret ของ LINE channel
// output: true ถ้า signature ตรงกับ expected signature
func verifySignature(body []byte, signature, channelSecret string) bool {
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
