package repositories

// normalizeLimit ปรับ limit ของ query list ให้อยู่ในช่วงที่ระบบยอมรับ
// input: limit จำนวนรายการจาก caller หรือ query string
// output: int limit โดย default เป็น 50 และสูงสุดไม่เกิน 200
func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}

// normalizeOffset ปรับ offset ของ query list ไม่ให้ติดลบ
// input: offset ตำแหน่งเริ่มต้นจาก caller หรือ query string
// output: int offset โดยค่าติดลบจะถูกแปลงเป็น 0
func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
