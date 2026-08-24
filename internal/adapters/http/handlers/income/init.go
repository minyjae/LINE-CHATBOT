package income

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type IIncomeController interface {
	List(c *fiber.Ctx) error
	Summary(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type incomeController struct {
	utilControllers.GenericController
	incomeService servicePort.IncomeService
}

func NewIncomeController(incomeService servicePort.IncomeService) IIncomeController {
	return &incomeController{
		GenericController: utilControllers.NewGenericController(),
		incomeService:     incomeService,
	}
}
