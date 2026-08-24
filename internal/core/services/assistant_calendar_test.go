package services

import (
	"strings"
	"testing"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

func TestIsCalendarListCommand(t *testing.T) {
	phrases := []string{
		"ดูนัดทั้งหมด",
		"ดู calendar ทั้งหมด",
		"มีนัดอะไรบ้าง",
		"เช็คนัดทั้งหมด",
		"รายการนัดหมาย",
	}

	for _, phrase := range phrases {
		t.Run(phrase, func(t *testing.T) {
			if !isCalendarListCommand(phrase) {
				t.Fatalf("expected phrase %q to be calendar list command", phrase)
			}
		})
	}
}

func TestFormatCalendarListReply(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	events := []*entities.CalendarEvent{{
		Title:   "ประชุมทีม",
		StartAt: time.Date(2026, time.August, 24, 10, 0, 0, 0, loc),
	}}

	reply := formatCalendarListReply(events, loc)
	if !strings.Contains(reply, "รายการนัดหมายทั้งหมด") {
		t.Fatalf("reply missing header: %q", reply)
	}
	if !strings.Contains(reply, "ประชุมทีม") {
		t.Fatalf("reply missing event title: %q", reply)
	}
	if !strings.Contains(reply, "24 Aug 2026 10:00") {
		t.Fatalf("reply missing event time: %q", reply)
	}
}
