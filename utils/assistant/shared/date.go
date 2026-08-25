package shared

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseReferenceDate แปลง phrase วันที่ในข้อความให้เป็นวันที่จริง
// input: text ข้อความดิบจากผู้ใช้, current เวลาปัจจุบัน, loc timezone ที่ต้องการใช้
// output: time.Time ที่ set เวลาเป็น 00:00:00 ใน loc, bool เป็น true เมื่อ parse วันที่ได้
func ParseReferenceDate(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	lower := strings.ToLower(text)
	if strings.Contains(text, "วันนี้") || strings.Contains(lower, "today") {
		return time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc), true
	}
	if strings.Contains(text, "เมื่อวาน") || strings.Contains(lower, "yesterday") {
		date := current.AddDate(0, 0, -1)
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), true
	}
	if strings.Contains(text, "พรุ่งนี้") || strings.Contains(lower, "tomorrow") {
		date := current.AddDate(0, 0, 1)
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), true
	}

	if date, ok := parseISODateInText(text, loc); ok {
		return date, true
	}
	if date, ok := parseSlashDateInText(text, current, loc); ok {
		return date, true
	}
	if date, ok := parseThaiNumericDatePhraseInText(text, current, loc); ok {
		return date, true
	}
	if date, ok := parseThaiMonthDateInText(text, current, loc); ok {
		return date, true
	}
	if date, ok := parseDayOnlyInText(text, current, loc); ok {
		return date, true
	}
	return time.Time{}, false
}

// parseISODateInText parse วันที่รูปแบบ YYYY-MM-DD จากข้อความ
// input: text ข้อความดิบจากผู้ใช้, loc timezone ที่ต้องการใช้
// output: time.Time วันที่ที่ parse ได้, bool เป็น true เมื่อเจอวันที่ถูกต้อง
func parseISODateInText(text string, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{4})-(\d{1,2})-(\d{1,2})\b`)
	match := re.FindStringSubmatch(text)
	if len(match) != 4 {
		return time.Time{}, false
	}
	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	day, _ := strconv.Atoi(match[3])
	return buildDate(year, time.Month(month), day, loc)
}

// parseSlashDateInText parse วันที่รูปแบบ D/M, D-M, D/M/YYYY หรือ D-M-YYYY
// input: text ข้อความดิบจากผู้ใช้, current ใช้เติมปีเมื่อผู้ใช้ไม่ระบุปี, loc timezone ที่ต้องการใช้
// output: time.Time วันที่ที่ parse ได้, bool เป็น true เมื่อเจอวันที่ถูกต้อง
func parseSlashDateInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{1,2})[/-](\d{1,2})(?:[/-](\d{2,4}))?\b`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	year := current.Year()
	if len(match) >= 4 && match[3] != "" {
		year, _ = strconv.Atoi(match[3])
		year = normalizeYear(year)
	}
	return buildDate(year, time.Month(month), day, loc)
}

// parseThaiNumericDatePhraseInText parse วันที่ภาษาไทยที่ระบุเลขวัน/เดือน/ปี
// input: text เช่น "วันที่ 7 เดือน 9 ปี 2026", current ใช้เติมปีเมื่อไม่ระบุปี, loc timezone ที่ต้องการใช้
// output: time.Time วันที่ที่ parse ได้, bool เป็น true เมื่อเจอวันที่ถูกต้อง
func parseThaiNumericDatePhraseInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`วันที่\s*(\d{1,2})\s*เดือน\s*(\d{1,2})(?:\s*ปี\s*(\d{2,4}))?`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	year := current.Year()
	if len(match) >= 4 && match[3] != "" {
		year, _ = strconv.Atoi(match[3])
		year = normalizeYear(year)
	}
	return buildDate(year, time.Month(month), day, loc)
}

// parseThaiMonthDateInText parse วันที่ที่ใช้ชื่อเดือนภาษาไทยหรืออังกฤษ
// input: text เช่น "7 กันยายน 2026" หรือ "7 Sep 2026", current ใช้เติมปีเมื่อไม่ระบุปี, loc timezone ที่ต้องการใช้
// output: time.Time วันที่ที่ parse ได้, bool เป็น true เมื่อเจอวันที่ถูกต้อง
func parseThaiMonthDateInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	monthPattern := ThaiMonthPattern()
	re := regexp.MustCompile(`(?i)(\d{1,2})\s*(` + monthPattern + `)(?:\s*(\d{2,4}))?`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	month, ok := thaiMonthNumber(match[2])
	if !ok {
		return time.Time{}, false
	}
	year := current.Year()
	if len(match) >= 4 && match[3] != "" {
		year, _ = strconv.Atoi(match[3])
		year = normalizeYear(year)
	}
	return buildDate(year, month, day, loc)
}

// parseDayOnlyInText parse phrase ที่มีเฉพาะ "วันที่ X" โดยใช้เดือน/ปีปัจจุบัน
// input: text เช่น "วันที่ 7", current ใช้เติมเดือนและปี, loc timezone ที่ต้องการใช้
// output: time.Time วันที่ที่ parse ได้, bool เป็น true เมื่อเจอวันที่ถูกต้อง
func parseDayOnlyInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`วันที่\s*(\d{1,2})`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	return buildDate(current.Year(), current.Month(), day, loc)
}

// ParseReportMonth แปลง phrase เดือนในข้อความให้เป็นวันแรกของเดือนนั้น
// input: text ข้อความดิบจากผู้ใช้, current ใช้เติมปีเมื่อไม่ระบุปี, loc timezone ที่ต้องการใช้
// output: time.Time วันแรกของเดือนเวลา 00:00:00, bool เป็น true เมื่อ parse เดือนได้
func ParseReportMonth(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	if start, ok := parseISOMonthInText(text, loc); ok {
		return start, true
	}

	monthPattern := ThaiMonthPattern()
	re := regexp.MustCompile(`(?i)(` + monthPattern + `)(?:\s*(\d{2,4}))?`)
	match := re.FindStringSubmatch(text)
	if len(match) >= 2 {
		month, ok := thaiMonthNumber(match[1])
		if ok {
			year := current.Year()
			if len(match) >= 3 && match[2] != "" {
				year, _ = strconv.Atoi(match[2])
				year = normalizeYear(year)
			}
			return time.Date(year, month, 1, 0, 0, 0, 0, loc), true
		}
	}

	numberMonth := regexp.MustCompile(`เดือน\s*(\d{1,2})(?:[/-](\d{2,4}))?`)
	match = numberMonth.FindStringSubmatch(text)
	if len(match) >= 2 {
		monthNumber, _ := strconv.Atoi(match[1])
		year := current.Year()
		if len(match) >= 3 && match[2] != "" {
			year, _ = strconv.Atoi(match[2])
			year = normalizeYear(year)
		}
		if monthNumber >= 1 && monthNumber <= 12 {
			return time.Date(year, time.Month(monthNumber), 1, 0, 0, 0, 0, loc), true
		}
	}

	return time.Time{}, false
}

// parseISOMonthInText parse เดือนรูปแบบ YYYY-MM จากข้อความ
// input: text ข้อความดิบจากผู้ใช้, loc timezone ที่ต้องการใช้
// output: time.Time วันแรกของเดือน, bool เป็น true เมื่อเดือนอยู่ในช่วง 1-12
func parseISOMonthInText(text string, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{4})-(\d{1,2})\b`)
	match := re.FindStringSubmatch(text)
	if len(match) != 3 {
		return time.Time{}, false
	}
	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	if month < 1 || month > 12 {
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc), true
}

// buildDate สร้าง time.Time พร้อม validate ว่าวัน/เดือน/ปีมีอยู่จริง
// input: year ปี ค.ศ. หรือ พ.ศ., month เดือน, day วันที่, loc timezone ที่ต้องการใช้
// output: time.Time เวลา 00:00:00, bool เป็น false ถ้าวันที่ invalid เช่น 31 กุมภาพันธ์
func buildDate(year int, month time.Month, day int, loc *time.Location) (time.Time, bool) {
	year = normalizeYear(year)
	if month < time.January || month > time.December || day < 1 || day > 31 {
		return time.Time{}, false
	}
	date := time.Date(year, month, day, 0, 0, 0, 0, loc)
	if date.Year() != year || date.Month() != month || date.Day() != day {
		return time.Time{}, false
	}
	return date, true
}

// normalizeYear แปลงปีที่ผู้ใช้พิมพ์ให้เป็น ค.ศ.
// input: year ปีแบบ 2 หลัก, ค.ศ. หรือ พ.ศ.
// output: int ปี ค.ศ. เช่น 26 เป็น 2026 และ 2569 เป็น 2026
func normalizeYear(year int) int {
	if year < 100 {
		return 2000 + year
	}
	if year > 2400 {
		return year - 543
	}
	return year
}

// ContainsThaiMonthName ตรวจว่าข้อความมีชื่อเดือนหรือคำย่อเดือนที่ระบบรู้จักหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบชื่อเดือนภาษาไทย/อังกฤษหรือคำย่อเดือน
func ContainsThaiMonthName(text string) bool {
	for name := range thaiMonthMap() {
		if strings.Contains(strings.ToLower(text), strings.ToLower(name)) {
			return true
		}
	}
	return false
}

// ThaiMonthPattern สร้าง regex pattern สำหรับจับชื่อเดือนและคำย่อเดือน
// input: ไม่มี
// output: string pattern ที่นำไปประกอบ regexp ได้
func ThaiMonthPattern() string {
	return strings.Join([]string{
		"มกราคม", "ม\\.ค\\.?", "jan", "january",
		"กุมภาพันธ์", "ก\\.พ\\.?", "feb", "february",
		"มีนาคม", "มี\\.ค\\.?", "mar", "march",
		"เมษายน", "เม\\.ย\\.?", "apr", "april",
		"พฤษภาคม", "พ\\.ค\\.?", "may",
		"มิถุนายน", "มิ\\.ย\\.?", "jun", "june",
		"กรกฎาคม", "ก\\.ค\\.?", "jul", "july",
		"สิงหาคม", "ส\\.ค\\.?", "aug", "august",
		"กันยายน", "ก\\.ย\\.?", "sep", "september",
		"ตุลาคม", "ต\\.ค\\.?", "oct", "october",
		"พฤศจิกายน", "พ\\.ย\\.?", "nov", "november",
		"ธันวาคม", "ธ\\.ค\\.?", "dec", "december",
	}, "|")
}

// thaiMonthNumber แปลงชื่อเดือนหรือคำย่อให้เป็น time.Month
// input: value ชื่อเดือน เช่น "กันยายน", "ก.ย.", "sep"
// output: time.Month ของเดือนนั้น, bool เป็น true เมื่อรู้จักชื่อเดือน
func thaiMonthNumber(value string) (time.Month, bool) {
	normalized := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(value, ".")))
	normalized = strings.ReplaceAll(normalized, ".", "")
	month, ok := thaiMonthMap()[normalized]
	return month, ok
}

// thaiMonthMap คืน mapping ชื่อเดือน/คำย่อ ไปยัง time.Month
// input: ไม่มี
// output: map[string]time.Month สำหรับใช้ lookup เดือน
func thaiMonthMap() map[string]time.Month {
	return map[string]time.Month{
		"มกราคม": time.January, "มค": time.January, "jan": time.January, "january": time.January,
		"กุมภาพันธ์": time.February, "กพ": time.February, "feb": time.February, "february": time.February,
		"มีนาคม": time.March, "มีค": time.March, "mar": time.March, "march": time.March,
		"เมษายน": time.April, "เมย": time.April, "apr": time.April, "april": time.April,
		"พฤษภาคม": time.May, "พค": time.May, "may": time.May,
		"มิถุนายน": time.June, "มิย": time.June, "jun": time.June, "june": time.June,
		"กรกฎาคม": time.July, "กค": time.July, "jul": time.July, "july": time.July,
		"สิงหาคม": time.August, "สค": time.August, "aug": time.August, "august": time.August,
		"กันยายน": time.September, "กย": time.September, "sep": time.September, "september": time.September,
		"ตุลาคม": time.October, "ตค": time.October, "oct": time.October, "october": time.October,
		"พฤศจิกายน": time.November, "พย": time.November, "nov": time.November, "november": time.November,
		"ธันวาคม": time.December, "ธค": time.December, "dec": time.December, "december": time.December,
	}
}

// StartOfReportWeek หาวันจันทร์ของสัปดาห์จากวันที่ที่ส่งเข้ามา
// input: date วันที่อ้างอิง, loc timezone ที่ต้องการใช้
// output: time.Time วันจันทร์เวลา 00:00:00 ของสัปดาห์นั้น
func StartOfReportWeek(date time.Time, loc *time.Location) time.Time {
	current := date.In(loc)
	daysSinceMonday := int(current.Weekday()) - int(time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7
	}
	startDate := current.AddDate(0, 0, -daysSinceMonday)
	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
}

// FormatReportPeriodLabel สร้าง label ภาษาไทยสำหรับช่วงรายงาน
// input: prefix ประเภทช่วงเวลา เช่น "วันที่", "สัปดาห์", "เดือน"; start/end ช่วงเวลา; loc timezone ที่ใช้แสดงผล
// output: string label ที่นำไปใส่ reply text ได้
func FormatReportPeriodLabel(prefix string, start, end time.Time, loc *time.Location) string {
	switch prefix {
	case "วันที่":
		return "วันที่ " + start.In(loc).Format("02 Jan 2006")
	case "สัปดาห์":
		return fmt.Sprintf("สัปดาห์ %s-%s", start.In(loc).Format("02 Jan"), end.Add(-time.Nanosecond).In(loc).Format("02 Jan 2006"))
	case "เดือน":
		return "เดือน " + start.In(loc).Format("Jan 2006")
	default:
		return prefix
	}
}
