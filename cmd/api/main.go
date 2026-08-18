package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"minyjae/go-starter/internal/adapters/http/routes"
	"minyjae/go-starter/internal/config"
	"minyjae/go-starter/utils"
)

// main: entry point ของ application
//
// ขั้นตอนการ bootstrap (เรียงตามลำดับ):
//   1. โหลด config + เปิด connection ไป database
//   2. wire-up repository → service ของแต่ละ domain
//      ตัวอย่าง (เพิ่มเองหลังสร้าง implementation):
//        exampleRepo := repositories.NewExampleRepository(db)
//        exampleService := services.NewExampleServiceImpl(exampleRepo)
//   3. สร้าง fiber app + register middleware (logger, cors, ...)
//   4. ลงทะเบียน route ผ่าน routes.SetupRoute โดยส่ง service เข้าไปเป็น dep
//   5. start server ที่ port ตาม config
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db := config.SetupDatabase(cfg)
	_ = db // TODO: ส่งต่อ db ไปยัง repository constructor ของแต่ละ domain

	// === wire-up service ของแต่ละ domain ที่นี่ ===
	// ดูตัวอย่างที่ internal/core/services/example.go และ
	//                 internal/adapters/presistance/repositories/example.go

	resp := utils.NewResponse()
	app := fiber.New(fiber.Config{
		ErrorHandler: resp.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsAllows,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	routes.SetupRoute(app)
	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
