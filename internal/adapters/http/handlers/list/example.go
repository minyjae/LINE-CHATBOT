package list

// ดู pattern เต็มได้ที่ internal/adapters/http/handlers/auth/example.go
//
// โครงสร้างไฟล์ใน folder นี้:
//   init.go        : interface + struct + constructor (NewListController)
//   <action>.go    : 1 ไฟล์ = 1 handler action (CreateList, GetLists, ...)
//
// ตัวอย่าง:
//
//   package list
//
//   import "github.com/gofiber/fiber/v2"
//
//   func (h *listController) ExampleAction(c *fiber.Ctx) error {
//       userID, ok := c.Locals("user_id").(uint)
//       if !ok || userID == 0 {
//           return h.Response.Unauthorized(c, "Missing user context", "MISSING_USER_CONTEXT")
//       }
//       return h.Response.Item(c, "ok", nil)
//   }
