package finance

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// IFinanceController กำหนด HTTP handler ของ finance endpoint
// input: Fiber context จาก route ที่เรียก Summary
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type IFinanceController interface {
	Summary(c *fiber.Ctx) error
}

// financeController เก็บ dependency สำหรับสรุปการเงินรวมรายรับ/รายจ่าย
// input: สร้างจาก NewFinanceController พร้อม ExpenseService และ IncomeService
// output: controller ที่รวม summary จากสอง service เป็น response เดียว
type financeController struct {
	utilControllers.GenericController
	expenseService servicePort.ExpenseService
	incomeService  servicePort.IncomeService
}

// NewFinanceController สร้าง finance controller
// input: expenseService และ incomeService สำหรับคำนวณ summary
// output: IFinanceController ที่พร้อมผูก route
func NewFinanceController(expenseService servicePort.ExpenseService, incomeService servicePort.IncomeService) IFinanceController {
	return &financeController{
		GenericController: utilControllers.NewGenericController(),
		expenseService:    expenseService,
		incomeService:     incomeService,
	}
}
