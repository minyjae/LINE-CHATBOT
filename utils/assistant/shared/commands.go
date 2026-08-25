package shared

import "strings"

// IsQuestion ตรวจว่า input text มีลักษณะเป็นคำถามหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบคำถาม เช่น "อะไร", "บ้าง", "ไหม", "เท่าไหร่" หรือเครื่องหมายคำถาม
func IsQuestion(text string) bool {
	questionWords := []string{"อะไร", "บ้าง", "ไหม", "มั้ย", "เท่าไหร่", "เท่าไร", "หรือยัง", "ยังไง", "?", "？"}
	return ContainsAny(text, questionWords)
}

// ContainsAny ตรวจว่า input text มี keyword ใด keyword หนึ่งอยู่หรือไม่ โดยไม่สนตัวพิมพ์เล็ก/ใหญ่
// input: text ข้อความดิบจากผู้ใช้, keywords ชุดคำที่ต้องการค้นหา
// output: true ถ้าพบอย่างน้อย 1 keyword ใน text
func ContainsAny(text string, keywords []string) bool {
	lower := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// HasPrefixAny ตรวจว่า input text ขึ้นต้นด้วย prefix ใด prefix หนึ่งหรือไม่ โดยไม่สนตัวพิมพ์เล็ก/ใหญ่
// input: text ข้อความดิบจากผู้ใช้, prefixes ชุดคำขึ้นต้นที่ต้องการตรวจ
// output: true ถ้า text เริ่มต้นด้วย prefix อย่างน้อย 1 ค่า
func HasPrefixAny(text string, prefixes []string) bool {
	lower := strings.ToLower(text)
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}
