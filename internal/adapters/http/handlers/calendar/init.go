package calendar

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type ICalendarController interface {
	List(c *fiber.Ctx) error
	ListByDate(c *fiber.Ctx) error
}

type calendarController struct {
	utilControllers.GenericController
	calendarService servicePort.CalendarEventService
}

func NewCalendarController(calendarService servicePort.CalendarEventService) ICalendarController {
	return &calendarController{
		GenericController: utilControllers.NewGenericController(),
		calendarService:   calendarService,
	}
}
