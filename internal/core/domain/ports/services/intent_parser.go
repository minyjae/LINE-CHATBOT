package services

import "minyjae/go-starter/types"

// AssistantIntentParser คือ contract สำหรับ parser ภายนอกหรือ AI ที่ตีความข้อความผู้ใช้
// input: IntentParseInput ที่มี text, now, timezone และ locale
// output: ParsedAssistantIntent ที่มี intent/entities/confidence หรือ error
type AssistantIntentParser interface {
	Parse(input types.IntentParseInput) (*types.ParsedAssistantIntent, error)
}
