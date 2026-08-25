package todo

import (
	"strings"

	shared "minyjae/go-starter/utils/assistant/shared"
)

// IsCommand ตรวจว่า text เป็นคำสั่งสร้าง todo หรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบคำว่า "todo", "ทูดู" หรือ "เพิ่มงาน"
func IsCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "todo") || strings.Contains(text, "ทูดู") || strings.Contains(text, "เพิ่มงาน")
}

// CleanupTitle ตัดคำสั่ง todo ออกจากข้อความ เพื่อเหลือชื่องาน
// input: text เช่น "เพิ่มงานซื้อของ"
// output: string title เช่น "ซื้อของ"
func CleanupTitle(text string) string {
	return shared.CleanupByRemoving(text, []string{"เพิ่ม todo", "เพิ่ม Todo", "todo", "Todo", "เพิ่มทูดู", "ทูดู", "เพิ่มงาน"})
}
