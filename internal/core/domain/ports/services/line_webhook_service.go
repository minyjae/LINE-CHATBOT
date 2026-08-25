package services

import "minyjae/go-starter/types"

// LineWebhookService คือ contract สำหรับประมวลผลข้อความ text ที่มาจาก LINE
// input: LineTextMessageInput ที่มี line user id, reply token, text และเวลา
// output: LineTextMessageResult ที่มี user/message log/intent/reply หรือ error
type LineWebhookService interface {
	HandleTextMessage(input types.LineTextMessageInput) (*types.LineTextMessageResult, error)
}
