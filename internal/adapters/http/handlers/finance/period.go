package finance

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SummaryPeriod คือช่วงเวลาที่ controller ใช้ query สรุปการเงิน
// input: สร้างจาก ParseSummaryPeriod จาก query period/date/month
// output: struct ที่มี kind/label/start/end สำหรับส่งต่อ service
type SummaryPeriod struct {
	Kind  string    `json:"kind"`
	Label string    `json:"label"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ParseSummaryPeriod แปลง query period/date/month ให้เป็นช่วงเวลา start/end
// input: periodValue เช่น day/week/month, dateValue YYYY-MM-DD, monthValue YYYY-MM, now, loc timezone
// output: SummaryPeriod หรือ error ถ้า period/date/month ไม่ถูกต้อง
func ParseSummaryPeriod(periodValue, dateValue, monthValue string, now time.Time, loc *time.Location) (SummaryPeriod, error) {
	if loc == nil {
		loc = time.Local
	}
	now = now.In(loc)

	period := strings.ToLower(strings.TrimSpace(periodValue))
	if period == "" {
		if strings.TrimSpace(dateValue) != "" {
			period = "day"
		} else {
			period = "month"
		}
	}

	date, hasDate, err := parseDateValue(dateValue, loc)
	if err != nil {
		return SummaryPeriod{}, err
	}
	if !hasDate {
		date = now
	}

	switch period {
	case "day", "daily":
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		return SummaryPeriod{Kind: "day", Label: "วันนี้", Start: start, End: start.AddDate(0, 0, 1)}, nil
	case "week", "weekly":
		start := startOfWeek(date, loc)
		return SummaryPeriod{Kind: "week", Label: "สัปดาห์นี้", Start: start, End: start.AddDate(0, 0, 7)}, nil
	case "month", "monthly":
		if strings.TrimSpace(monthValue) != "" {
			year, month, err := parseMonthValue(monthValue)
			if err != nil {
				return SummaryPeriod{}, err
			}
			date = time.Date(year, month, 1, 0, 0, 0, 0, loc)
		}
		start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, loc)
		return SummaryPeriod{Kind: "month", Label: "เดือนนี้", Start: start, End: start.AddDate(0, 1, 0)}, nil
	default:
		return SummaryPeriod{}, fiber.ErrBadRequest
	}
}

// parseDateValue แปลง query date เป็น time.Time
// input: value วันที่รูปแบบ YYYY-MM-DD หรือค่าว่าง, loc timezone
// output: time.Time, bool ว่ามี date ส่งมาไหม, และ error ถ้า format ไม่ถูกต้อง
func parseDateValue(value string, loc *time.Location) (time.Time, bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false, nil
	}
	date, err := time.ParseInLocation("2006-01-02", value, loc)
	if err != nil {
		return time.Time{}, false, err
	}
	return date, true, nil
}

// parseMonthValue แปลง query month เป็นปีและเดือน
// input: value เดือนรูปแบบ YYYY-MM
// output: year, time.Month และ error ถ้า format หรือเลขเดือนไม่ถูกต้อง
func parseMonthValue(value string) (int, time.Month, error) {
	parts := strings.Split(strings.TrimSpace(value), "-")
	if len(parts) != 2 {
		return 0, 0, fiber.ErrBadRequest
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	monthNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if monthNumber < 1 || monthNumber > 12 {
		return 0, 0, fiber.ErrBadRequest
	}
	return year, time.Month(monthNumber), nil
}

// startOfWeek หาวันจันทร์ของสัปดาห์จากวันที่ที่ส่งมา
// input: date วันที่อ้างอิง, loc timezone
// output: time.Time วันจันทร์เวลา 00:00:00 ของสัปดาห์นั้น
func startOfWeek(date time.Time, loc *time.Location) time.Time {
	current := date.In(loc)
	daysSinceMonday := int(current.Weekday()) - int(time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7
	}
	startDate := current.AddDate(0, 0, -daysSinceMonday)
	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
}
