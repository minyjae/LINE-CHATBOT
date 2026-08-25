package calendar

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// ICalendarController กำหนด HTTP handler ทั้งหมดของ calendar endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type ICalendarController interface {
	List(c *fiber.Ctx) error
	ListByDate(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// calendarController เก็บ dependency สำหรับ calendar HTTP handler
// input: สร้างจาก NewCalendarController พร้อม CalendarEventService
// output: controller ที่แปลง HTTP request เป็น service call
type calendarController struct {
	utilControllers.GenericController
	calendarService servicePort.CalendarEventService
}

// NewCalendarController สร้าง calendar controller
// input: calendarService สำหรับทำ business logic ของ calendar event
// output: ICalendarController ที่พร้อมผูก route
func NewCalendarController(calendarService servicePort.CalendarEventService) ICalendarController {
	return &calendarController{
		GenericController: utilControllers.NewGenericController(),
		calendarService:   calendarService,
	}
}
