package reminder

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type IReminderController interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type reminderController struct {
	utilControllers.GenericController
	reminderService servicePort.ReminderService
}

func NewReminderController(reminderService servicePort.ReminderService) IReminderController {
	return &reminderController{
		GenericController: utilControllers.NewGenericController(),
		reminderService:   reminderService,
	}
}
