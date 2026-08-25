package services

import "minyjae/go-starter/types"

// AssistantService คือ contract สำหรับตีความข้อความผู้ใช้และคืนผลลัพธ์ของ assistant
// input: AssistantMessageInput ที่มี user, message log, text, เวลา และ timezone
// output: AssistantMessageResult ที่มี intent/reply/data หรือ error
type AssistantService interface {
	HandleTextMessage(input types.AssistantMessageInput) (*types.AssistantMessageResult, error)
}
