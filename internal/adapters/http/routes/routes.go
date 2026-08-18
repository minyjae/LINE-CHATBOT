package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoute: รวมการลงทะเบียน route ของทุก domain ไว้ในที่เดียว
//
// วิธีเพิ่ม route ใหม่:
//   1. สร้างไฟล์ <domain>.go ใน folder นี้ (ดู example.go)
//   2. ใน <domain>.go ประกาศฟังก์ชัน <Domain>Route(app, deps...)
//   3. เพิ่ม service ที่ต้องใช้เป็นพารามิเตอร์ของ SetupRoute แล้วเรียกฟังก์ชันนั้นใน body
//
// ตัวอย่าง signature (เพิ่ม service เข้าไปเองตามที่ใช้):
//
//   func SetupRoute(
//       app *fiber.App,
//       exampleService servicePort.ExampleService,
//   ) {
//       ExampleRoute(app, exampleService)
//   }
func SetupRoute(app *fiber.App) {
	// ลงทะเบียน route ของแต่ละ domain ที่นี่
}
