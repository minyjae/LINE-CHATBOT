package shared

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HasExplicitTime ตรวจว่าข้อความมีเวลาที่ผู้ใช้ระบุชัดเจนหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้า ExtractHourMinute parse ชั่วโมง/นาทีได้
func HasExplicitTime(text string) bool {
	_, _, ok := ExtractHourMinute(text)
	return ok
}

// ParseNaturalTime แปลงข้อความวันที่/เวลาให้เป็น time.Time เดียวกัน
// input: text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone, defaultHour/defaultMinute เวลาตั้งต้นเมื่อมีวันที่แต่ไม่มีเวลา
// output: time.Time วันที่/เวลาที่ parse ได้, bool เป็น true เมื่อพบวันที่หรือเวลาจากข้อความ
func ParseNaturalTime(text string, now time.Time, loc *time.Location, defaultHour, defaultMinute int) (time.Time, bool) {
	date := now.In(loc)
	hasDate := false
	if parsedDate, ok := ParseReferenceDate(text, date, loc); ok {
		date = parsedDate
		hasDate = true
	} else if strings.Contains(text, "พรุ่งนี้") {
		date = date.AddDate(0, 0, 1)
		hasDate = true
	}

	hour, minute, hasTime := ExtractHourMinute(text)
	if !hasTime {
		hour = defaultHour
		minute = defaultMinute
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc), hasDate || hasTime
}

// ExtractHourMinute ดึงชั่วโมงและนาทีจากข้อความเวลาแบบไทย/ตัวเลข
// input: text ข้อความดิบ เช่น "9 โมง 40", "09:40", "บ่าย 2", "ตี 5", "เที่ยงครึ่ง"
// output: hour, minute ในรูปแบบ 24 ชั่วโมง และ bool เป็น true เมื่อ parse ได้
func ExtractHourMinute(text string) (int, int, bool) {
	if hour, minute, ok := extractThaiNoonTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiTuumTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiAfternoonTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiClockTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractSeparatedHourMinute(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiShortClockTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiDawnTime(text); ok {
		return hour, minute, true
	}

	return 0, 0, false
}

// extractThaiNoonTime parse เวลาแบบ "เที่ยง" หรือ "เที่ยงครึ่ง"
// input: text ข้อความดิบจากผู้ใช้
// output: hour/minute เช่น 12:00 หรือ 12:30, bool เป็น true เมื่อ match
func extractThaiNoonTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`เที่ยง(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?`)
	if match := re.FindStringSubmatch(text); len(match) == 2 {
		minute := parseThaiMinuteSuffix(match[1])
		return validHourMinute(12, minute)
	}
	return 0, 0, false
}

// extractSeparatedHourMinute parse เวลาแบบมีตัวคั่น เช่น "09:40" หรือ "9.40"
// input: text ข้อความดิบจากผู้ใช้
// output: hour/minute หลัง normalize เป็นเวลา 24 ชั่วโมง, bool เป็น true เมื่อ match
func extractSeparatedHourMinute(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})[:.](\d{2})\s*(?:น\.?|นาฬิกา)?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])
		return validHourMinute(normalizeThaiHour(text, hour), minute)
	}
	return 0, 0, false
}

// extractThaiClockTime parse เวลาแบบ "X โมง" พร้อมคำบอกช่วงวัน
// input: text เช่น "9 โมงเช้า", "6 โมงเย็น", "9 โมง 40 นาที"
// output: hour/minute หลัง normalize เป็นเวลา 24 ชั่วโมง, bool เป็น true เมื่อ match
func extractThaiClockTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*โมง(?:\s*(เช้า|เย็น|ค่ำ))?(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?(?:\s*(เช้า|เย็น|ค่ำ))?`)
	if match := re.FindStringSubmatch(text); len(match) == 5 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(firstNonEmpty(match[3]))
		context := strings.TrimSpace(match[0] + " " + match[2] + " " + match[4])
		return validHourMinute(normalizeThaiHour(context, hour), minute)
	}
	return 0, 0, false
}

// extractThaiShortClockTime parse เวลาแบบสั้นที่ลงท้ายด้วย "น." หรือ "นาฬิกา"
// input: text เช่น "9 น.", "9 นาฬิกา", "9 น.ครึ่ง"
// output: hour/minute หลัง normalize เป็นเวลา 24 ชั่วโมง, bool เป็น true เมื่อ match
func extractThaiShortClockTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*(?:น\.|นาฬิกา)(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(normalizeThaiHour(match[0], hour), minute)
	}
	return 0, 0, false
}

// extractThaiAfternoonTime parse เวลาแบบ "บ่าย X" หรือ "บ่ายโมง"
// input: text เช่น "บ่ายโมง", "บ่ายโมงครึ่ง", "บ่าย 2", "บ่าย 2 ครึ่ง", "บ่าย 2 โมง 15 นาที"
// output: hour/minute ในรูปแบบ 24 ชั่วโมง, bool เป็น true เมื่อ match
func extractThaiAfternoonTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`บ่าย\s*(?:โมง|(\d{1,2})(?:\s*โมง)?)(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour := 1
		if match[1] != "" {
			hour, _ = strconv.Atoi(match[1])
		}
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(normalizeThaiHour(match[0], hour), minute)
	}
	return 0, 0, false
}

// extractThaiTuumTime parse เวลาแบบ "X ทุ่ม"
// input: text เช่น "3 ทุ่ม", "3 ทุ่มครึ่ง"
// output: hour/minute ในรูปแบบ 24 ชั่วโมง เช่น 21:00, bool เป็น true เมื่อ match
func extractThaiTuumTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*ทุ่ม(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		tuum, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		hour := 18 + tuum
		if hour == 24 {
			hour = 0
		}
		return validHourMinute(hour, minute)
	}
	return 0, 0, false
}

// extractThaiDawnTime parse เวลาแบบ "ตี X"
// input: text เช่น "ตี 5", "ตี 5 ครึ่ง"
// output: hour/minute ในช่วงเช้ามืด, bool เป็น true เมื่อ match
func extractThaiDawnTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`ตี\s*(\d{1,2})(?:\s*(ครึ่ง|\d{1,2}(?:\s*นาที)?))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(hour, minute)
	}
	return 0, 0, false
}

// parseThaiMinuteSuffix แปลง suffix นาทีภาษาไทยเป็นตัวเลขนาที
// input: value เช่น "", "ครึ่ง", "40 นาที"
// output: int นาที โดย "" เป็น 0 และ "ครึ่ง" เป็น 30
func parseThaiMinuteSuffix(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if strings.Contains(value, "ครึ่ง") {
		return 30
	}
	re := regexp.MustCompile(`(\d{1,2})`)
	if match := re.FindStringSubmatch(value); len(match) == 2 {
		minute, _ := strconv.Atoi(match[1])
		return minute
	}
	return 0
}

// normalizeThaiHour ปรับชั่วโมงตามบริบทภาษาไทยให้เป็นเวลา 24 ชั่วโมง
// input: text บริบทเวลา, hour ชั่วโมงที่ parse ได้
// output: int ชั่วโมงหลังปรับ เช่น "บ่าย 2" เป็น 14, "6 โมงเย็น" เป็น 18
func normalizeThaiHour(text string, hour int) int {
	if strings.Contains(text, "เที่ยงคืน") {
		return 0
	}
	if strings.Contains(text, "เที่ยง") && hour == 12 {
		return 12
	}
	if strings.Contains(text, "บ่าย") && hour < 12 {
		return hour + 12
	}
	if (strings.Contains(text, "เย็น") || strings.Contains(text, "ค่ำ")) && hour < 12 {
		return hour + 12
	}
	return hour
}

// validHourMinute validate ขอบเขตชั่วโมงและนาที
// input: hour ชั่วโมง, minute นาที
// output: hour/minute เดิมพร้อม true ถ้าอยู่ในช่วงถูกต้อง, หรือ 0/0/false ถ้า invalid
func validHourMinute(hour, minute int) (int, int, bool) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, false
	}
	return hour, minute, true
}
