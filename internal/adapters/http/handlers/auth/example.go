package auth

// Folder: internal/adapters/http/handlers/<domain>
//
// หน้าที่: รวม HTTP handler ของ domain หนึ่งๆ (เช่น auth, user, list)
//
// กติกาตั้งชื่อไฟล์:
//   - init.go      : ประกาศ interface + struct + constructor (Wire-up dependency)
//   - <action>.go  : implement handler method แต่ละ action (1 ไฟล์ = 1 action)
//
// แต่ละ action ใช้ pointer receiver ของ struct ใน init.go เพื่อเข้าถึง
// service และ utility ที่ฝังไว้ผ่าน GenericController (Response, Pagination)
//
// ตัวอย่าง action (สมมุติเพิ่มไฟล์ logout.go):
//
//   package auth
//
//   import "github.com/gofiber/fiber/v2"
//
//   func (h *authController) Logout(c *fiber.Ctx) error {
//       // 1. ดึงข้อมูลจาก context (c.Locals หลัง JWT middleware)
//       // 2. เรียก service: h.authService.Logout(...)
//       // 3. ตอบกลับด้วย helper: h.Response.Item / Created / BadRequest ...
//       return h.Response.Item(c, "logout successfully", nil)
//   }
//
// อย่าลืม:
//   - เพิ่ม method ลงใน interface IAuthController ที่ init.go
//   - ผูก route ใหม่ใน internal/adapters/http/routes/<domain>.go
