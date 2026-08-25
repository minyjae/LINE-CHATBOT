package calendar

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
			if !IsListCommand(phrase) {
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

	reply := FormatListReply(events, loc)
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

func TestCleanupCalendarTitleCutsTimeClause(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{text: "นัดดูรถ ตอน บ่าย2 30 นาที", want: "ดูรถ"},
		{text: "นัดดูรถที่ โตโยต้านครพิงค์ ตอนเที่ยง 30 นาที", want: "ดูรถที่ โตโยต้านครพิงค์"},
		{text: "ลงตาราง ดูรถที่ โตโยต้า ตอนเที่ยงครึ่ง", want: "ดูรถที่ โตโยต้า"},
		{text: "นัดดูรถที่ โจเส่ เที่ยงครึ่ง", want: "ดูรถที่ โจเส่"},
		{text: "นัดประชุม 30/08/2026 10 โมง", want: "ประชุม"},
		{text: "นัดประชุม 30 สิงหาคม 2026 10 โมง", want: "ประชุม"},
		{text: "นัดประชุมวันที่ 30 10 โมง", want: "ประชุม"},
		{text: "นัดดูจองรถ วันที่ 7 เดือน 9 ปี 2026 ตอน9โมง40", want: "ดูจองรถ"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := CleanupTitle(tt.text)
			if got != tt.want {
				t.Fatalf("title = %q, want %q", got, tt.want)
			}
		})
	}
}
