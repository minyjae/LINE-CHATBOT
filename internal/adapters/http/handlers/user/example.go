package user

// ดู pattern เต็มได้ที่ internal/adapters/http/handlers/auth/example.go
//
// โครงสร้างไฟล์ใน folder นี้:
//   init.go        : interface IUserController + struct + constructor
//   <action>.go    : 1 ไฟล์ = 1 handler action (GetUser, UpdateUser, ...)
//
// ตัวอย่าง:
//
//   package user
//
//   import "github.com/gofiber/fiber/v2"
//
//   func (h *userController) ExampleAction(c *fiber.Ctx) error {
//       return h.Response.Item(c, "ok", nil)
//   }
