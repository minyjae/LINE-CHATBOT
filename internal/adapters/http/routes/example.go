package routes

// Folder: internal/adapters/http/routes
//
// หน้าที่: ผูก URL path → handler ของแต่ละ domain
//
// กติกาตั้งชื่อไฟล์:
//   routes.go      : รวม register function ของทุก domain ไว้ใน SetupRoute
//   <domain>.go    : 1 ไฟล์ = 1 domain (auth.go, user.go, list.go, ...)
//
// แต่ละไฟล์ domain จะ:
//   1. สร้าง handler ผ่าน constructor ของ handler package (เช่น auth.NewAuthController)
//   2. ประกาศ Group + ผูก method (Get/Post/Put/Patch/Delete)
//   3. ถ้า route ต้องการการยืนยันตัวตน ให้ใส่ middleware.JWTProtected()
//
// ตัวอย่างไฟล์ example_route.go:
//
//   package routes
//
//   import (
//       "github.com/gofiber/fiber/v2"
//       "minyjae/go-starter/internal/adapters/http/handlers/example"
//       "minyjae/go-starter/internal/adapters/http/middleware"
//       "minyjae/go-starter/internal/core/domain/ports/services"
//   )
//
//   func ExampleRoute(app *fiber.App, exampleService services.ExampleService) {
//       handler := example.NewExampleController(exampleService)
//
//       group := app.Group("/examples", middleware.JWTProtected())
//       group.Get("/", handler.List)
//       group.Post("/", handler.Create)
//   }
//
// จากนั้นเรียก ExampleRoute(...) ใน SetupRoute (ดู routes.go)
