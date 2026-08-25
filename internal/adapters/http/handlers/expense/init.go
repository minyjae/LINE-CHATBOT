package expense

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// IExpenseController กำหนด HTTP handler ทั้งหมดของ expense endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type IExpenseController interface {
	List(c *fiber.Ctx) error
	Summary(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// expenseController เก็บ dependency สำหรับ expense HTTP handler
// input: สร้างจาก NewExpenseController พร้อม ExpenseService
// output: controller ที่แปลง HTTP request เป็น service call
type expenseController struct {
	utilControllers.GenericController
	expenseService servicePort.ExpenseService
}

// NewExpenseController สร้าง expense controller
// input: expenseService สำหรับทำ business logic ของรายจ่าย
// output: IExpenseController ที่พร้อมผูก route
func NewExpenseController(expenseService servicePort.ExpenseService) IExpenseController {
	return &expenseController{
		GenericController: utilControllers.NewGenericController(),
		expenseService:    expenseService,
	}
}
