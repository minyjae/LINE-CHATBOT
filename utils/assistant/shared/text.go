package shared

import (
	"regexp"
	"strings"
)

// firstNonEmpty คืนค่าข้อความตัวแรกที่ไม่ว่างหลัง trim space
// input: values รายการข้อความที่ต้องการเลือก
// output: string ข้อความแรกที่ไม่ว่าง หรือ "" ถ้าทุกค่าเป็นค่าว่าง
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// NormalizeDescriptionSpaces จัดช่องว่างในข้อความให้สะอาดก่อนใช้เป็น title/description
// input: text ข้อความที่อาจมีช่องว่างซ้ำหรือช่องว่างหน้า punctuation
// output: string ข้อความที่ trim แล้ว ลดช่องว่างซ้ำเหลือ 1 ช่อง
func NormalizeDescriptionSpaces(text string) string {
	text = strings.TrimSpace(text)
	spaces := regexp.MustCompile(`\s+`)
	text = spaces.ReplaceAllString(text, " ")
	punctuationSpaces := regexp.MustCompile(`\s+([,.;:!?])`)
	return punctuationSpaces.ReplaceAllString(text, "$1")
}

// RemoveDatePhrases ลบ phrase วันที่ออกจากข้อความ เพื่อให้เหลือเฉพาะ title/description
// input: text ข้อความที่อาจมีวันที่ เช่น "2026-09-07", "7/9/2026", "วันที่ 7 เดือน 9"
// output: string ข้อความหลังตัดส่วนวันที่ออก
func RemoveDatePhrases(text string) string {
	patterns := []string{
		`\b\d{4}-\d{1,2}-\d{1,2}\b`,
		`\b\d{1,2}[/-]\d{1,2}(?:[/-]\d{2,4})?\b`,
		`วันที่\s*\d{1,2}\s*เดือน\s*\d{1,2}(?:\s*ปี\s*\d{2,4})?`,
		`วันที่\s*\d{1,2}\s*(?:` + ThaiMonthPattern() + `)?(?:\s*\d{2,4})?`,
		`\d{1,2}\s*(?:` + ThaiMonthPattern() + `)(?:\s*\d{2,4})?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

// RemoveTimePhrases ลบ phrase เวลาออกจากข้อความ เพื่อให้เหลือเฉพาะ title/description
// input: text ข้อความที่อาจมีเวลา เช่น "9 โมง", "09:40", "บ่าย 2", "ตี 5"
// output: string ข้อความหลังตัดส่วนเวลาออก
func RemoveTimePhrases(text string) string {
	patterns := []string{
		`เที่ยง(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`\d{1,2}\s*โมง(?:\s*(?:เช้า|เย็น|ค่ำ))?(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?(?:\s*(?:เช้า|เย็น|ค่ำ))?`,
		`\d{1,2}\s*ทุ่ม(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`บ่าย\s*\d{1,2}(?:\s*โมง)?(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`ตี\s*\d{1,2}(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`\d{1,2}[:.]\d{2}\s*(?:น\.?|นาฬิกา)?`,
		`\d{1,2}\s*(?:น\.|นาฬิกา)(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

// CleanupByRemoving ลบ token ที่ระบุออกจากข้อความแบบตรงตัว
// input: text ข้อความต้นฉบับ, tokens คำหรือวลีที่ต้องการลบ
// output: string ข้อความที่ถูกลบ token แล้วและ trim space
func CleanupByRemoving(text string, tokens []string) string {
	result := text
	for _, token := range tokens {
		result = strings.ReplaceAll(result, token, "")
	}
	return strings.TrimSpace(result)
}

// CleanupTimeWords ลบคำบอกเวลาแบบทั่วไปออกจากข้อความ
// input: text ข้อความที่อาจเหลือคำเวลา เช่น "วันนี้", "พรุ่งนี้", "ตอน", "เวลา" หรือเลขเวลา
// output: string ข้อความหลังลบคำเวลาเบื้องต้นและ trim space
func CleanupTimeWords(text string) string {
	replacers := []string{"วันนี้", "พรุ่งนี้", "โมง", "ตอน", "เวลา"}
	result := text
	for _, token := range replacers {
		result = strings.ReplaceAll(result, token, "")
	}
	re := regexp.MustCompile(`\d{1,2}(:\d{2})?`)
	result = re.ReplaceAllString(result, "")
	return strings.TrimSpace(result)
}
