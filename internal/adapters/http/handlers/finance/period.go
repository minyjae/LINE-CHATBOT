package finance

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SummaryPeriod struct {
	Kind  string    `json:"kind"`
	Label string    `json:"label"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

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

func startOfWeek(date time.Time, loc *time.Location) time.Time {
	current := date.In(loc)
	daysSinceMonday := int(current.Weekday()) - int(time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7
	}
	startDate := current.AddDate(0, 0, -daysSinceMonday)
	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
}
