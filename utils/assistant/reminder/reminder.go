package reminder

import (
	"strings"

	shared "minyjae/go-starter/utils/assistant/shared"
)

// IsCommand ตรวจว่า text เป็นคำสั่งสร้าง reminder หรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบคำว่า "เตือน"
func IsCommand(text string) bool {
	return strings.Contains(text, "เตือน")
}

// CleanupTitle ตัดคำสั่ง วันที่ และเวลาออกจาก reminder เพื่อเหลือ title ที่ต้องเตือน
// input: text เช่น "เตือนดูจองรถ ตอน 9 โมง 40 วันที่ 7 เดือน 9 ปี 2026"
// output: string title เช่น "ดูจองรถ"
func CleanupTitle(text string) string {
	cleaned := shared.CleanupByRemoving(text, []string{"เตือน", "remind me", "reminder"})
	cleaned = shared.RemoveDatePhrases(cleaned)
	cleaned = shared.RemoveTimePhrases(cleaned)
	cleaned = shared.CleanupTimeWords(cleaned)
	return shared.NormalizeDescriptionSpaces(cleaned)
}
