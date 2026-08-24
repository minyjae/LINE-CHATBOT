package finance

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type IFinanceController interface {
	Summary(c *fiber.Ctx) error
}

type financeController struct {
	utilControllers.GenericController
	expenseService servicePort.ExpenseService
	incomeService  servicePort.IncomeService
}

func NewFinanceController(expenseService servicePort.ExpenseService, incomeService servicePort.IncomeService) IFinanceController {
	return &financeController{
		GenericController: utilControllers.NewGenericController(),
		expenseService:    expenseService,
		incomeService:     incomeService,
	}
}
