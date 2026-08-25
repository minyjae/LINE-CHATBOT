package income

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// IIncomeController กำหนด HTTP handler ทั้งหมดของ income endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type IIncomeController interface {
	List(c *fiber.Ctx) error
	Summary(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// incomeController เก็บ dependency สำหรับ income HTTP handler
// input: สร้างจาก NewIncomeController พร้อม IncomeService
// output: controller ที่แปลง HTTP request เป็น service call
type incomeController struct {
	utilControllers.GenericController
	incomeService servicePort.IncomeService
}

// NewIncomeController สร้าง income controller
// input: incomeService สำหรับทำ business logic ของรายรับ
// output: IIncomeController ที่พร้อมผูก route
func NewIncomeController(incomeService servicePort.IncomeService) IIncomeController {
	return &incomeController{
		GenericController: utilControllers.NewGenericController(),
		incomeService:     incomeService,
	}
}
