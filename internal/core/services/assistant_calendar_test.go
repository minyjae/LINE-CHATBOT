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

func TestCleanupCalendarTitleCutsTimeClause(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{text: "นัดดูรถ ตอน บ่าย2 30 นาที", want: "ดูรถ"},
		{text: "นัดดูรถที่ โตโยต้านครพิงค์ ตอนเที่ยง 30 นาที", want: "ดูรถที่ โตโยต้านครพิงค์"},
		{text: "ลงตาราง ดูรถที่ โตโยต้า ตอนเที่ยงครึ่ง", want: "ดูรถที่ โตโยต้า"},
		{text: "นัดดูรถที่ โจเส่ เที่ยงครึ่ง", want: "ดูรถที่ โจเส่"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := cleanupCalendarTitle(tt.text)
			if got != tt.want {
				t.Fatalf("title = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseNaturalTimeWithAfternoonMinute(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	got, ok := parseNaturalTime("นัดดูรถ ตอน บ่าย2 30 นาที", now, loc, 9, 0)
	if !ok {
		t.Fatal("expected time")
	}
	want := time.Date(2026, time.August, 24, 14, 30, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestParseNaturalTimeWithThaiNoon(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	tests := []struct {
		text       string
		wantHour   int
		wantMinute int
	}{
		{text: "นัดดูรถที่ โตโยต้านครพิงค์ ตอนเที่ยง 30 นาที", wantHour: 12, wantMinute: 30},
		{text: "ลงตาราง ดูรถที่ โตโยต้า ตอนเที่ยงครึ่ง", wantHour: 12, wantMinute: 30},
		{text: "นัดดูรถที่ โจเส่ เที่ยงครึ่ง", wantHour: 12, wantMinute: 30},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got, ok := parseNaturalTime(tt.text, now, loc, 9, 0)
			if !ok {
				t.Fatal("expected time")
			}
			want := time.Date(2026, time.August, 24, tt.wantHour, tt.wantMinute, 0, 0, loc)
			if !got.Equal(want) {
				t.Fatalf("time = %s, want %s", got, want)
			}
		})
	}
}
