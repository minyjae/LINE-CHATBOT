package validators

// Folder: validators
//
// หน้าที่: ประกาศ request DTO + validate tag และ wrapper สำหรับเช็คผ่าน
//          go-playground/validator
//
// โครงสร้างไฟล์:
//   init.go        : utility ValidateStruct[T] + ErrorResponse (ใช้ร่วมทุก validator)
//   <domain>.go    : 1 ไฟล์ = 1 กลุ่ม request (user.go, list.go, ...)
//                    ภายในประกาศ struct ของ request + ฟังก์ชัน validator
//
// กติกาตั้งชื่อ:
//   - struct ลงท้ายด้วย Request (เช่น RegisterRequest)
//   - validator function ลงท้ายด้วย Validator (เช่น RegisterRequestValidator)
//     เพื่อเสียบเป็น middleware ก่อนถึง handler ได้
//
// ตัวอย่าง:
//
//   package validators
//
//   import "github.com/gofiber/fiber/v2"
//
//   type ExampleRequest struct {
//       Name string `json:"name" validate:"required"`
//   }
//
//   func ExampleRequestValidator(c *fiber.Ctx) error {
//       return ValidateStruct[ExampleRequest](c)
//   }
//
// การใช้งานใน route:
//   group.Post("/", validators.ExampleRequestValidator, handler.Create)
//
// ใน handler ดึง request body ที่ validate แล้วจาก c.Locals("request")
