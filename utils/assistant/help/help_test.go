package help

import (
	"strings"
	"testing"
)

func TestIsCommand(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "mina typo help", text: "เลขามินะช่วยต้วย", want: true},
		{name: "mina proper help", text: "เลขามินะช่วยด้วย", want: true},
		{name: "usage", text: "ใช้ยังไง", want: true},
		{name: "english", text: "help", want: true},
		{name: "not help", text: "เตือนทำการบ้านตอนบ่ายโมง", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCommand(tt.text); got != tt.want {
				t.Fatalf("IsCommand(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestFormatReplyIncludesFeatureExamples(t *testing.T) {
	reply := FormatReply()
	required := []string{
		"เตือนทำการบ้าน วันที่ 26 สิงหาคม 2026 ตอน 11 โมงครึ่ง",
		"ดูนัดทั้งหมด",
		"สรุปรายจ่ายเดือนนี้",
	}
	for _, phrase := range required {
		if !strings.Contains(reply, phrase) {
			t.Fatalf("reply does not include %q", phrase)
		}
	}
}
