package expense

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type IExpenseController interface {
	List(c *fiber.Ctx) error
	Summary(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type expenseController struct {
	utilControllers.GenericController
	expenseService servicePort.ExpenseService
}

func NewExpenseController(expenseService servicePort.ExpenseService) IExpenseController {
	return &expenseController{
		GenericController: utilControllers.NewGenericController(),
		expenseService:    expenseService,
	}
}
