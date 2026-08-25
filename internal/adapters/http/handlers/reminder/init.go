package reminder

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// IReminderController กำหนด HTTP handler ทั้งหมดของ reminder endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type IReminderController interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// reminderController เก็บ dependency สำหรับ reminder HTTP handler
// input: สร้างจาก NewReminderController พร้อม ReminderService
// output: controller ที่แปลง HTTP request เป็น service call
type reminderController struct {
	utilControllers.GenericController
	reminderService servicePort.ReminderService
}

// NewReminderController สร้าง reminder controller
// input: reminderService สำหรับทำ business logic ของ reminder
// output: IReminderController ที่พร้อมผูก route
func NewReminderController(reminderService servicePort.ReminderService) IReminderController {
	return &reminderController{
		GenericController: utilControllers.NewGenericController(),
		reminderService:   reminderService,
	}
}
