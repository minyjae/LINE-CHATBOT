package summary

import "strings"

// IsTomorrowCommand ตรวจว่า text เป็นคำถาม summary ของวันพรุ่งนี้หรือไม่
// input: text ข้อความดิบจากผู้ใช้ เช่น "พรุ่งนี้มีอะไรบ้าง"
// output: true ถ้ามีคำว่า "พรุ่งนี้" และถามว่า "มีอะไร" หรือ "ทำอะไร"
func IsTomorrowCommand(text string) bool {
	return strings.Contains(text, "พรุ่งนี้") && (strings.Contains(text, "มีอะไร") || strings.Contains(text, "ทำอะไร"))
}
