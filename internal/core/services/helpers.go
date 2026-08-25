package services

// defaultString คืนค่า fallback เมื่อ value ว่าง
// input: value ค่าที่ต้องการใช้, fallback ค่าสำรอง
// output: string เป็น value ถ้าไม่ว่าง ไม่อย่างนั้นเป็น fallback
func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
