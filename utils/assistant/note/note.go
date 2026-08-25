package note

import (
	"strings"

	shared "minyjae/go-starter/utils/assistant/shared"
)

// IsCommand ตรวจว่า text เป็นคำสั่งจด note หรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบคำว่า "จด" หรือ "note"
func IsCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(text, "จด") || strings.Contains(lower, "note")
}

// CleanupContent ตัดคำสั่ง note ออกจากข้อความ เพื่อเหลือเนื้อหาที่ต้องจด
// input: text เช่น "จดไว้ว่าโทรหาลูกค้า"
// output: string content เช่น "โทรหาลูกค้า"
func CleanupContent(text string) string {
	return shared.CleanupByRemoving(text, []string{"จดไว้ว่า", "จดว่า", "จด", "note", "Note"})
}
