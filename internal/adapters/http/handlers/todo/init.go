package todo

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// ITodoController กำหนด HTTP handler ทั้งหมดของ todo endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type ITodoController interface {
	List(c *fiber.Ctx) error
	ListPending(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// todoController เก็บ dependency สำหรับ todo HTTP handler
// input: สร้างจาก NewTodoController พร้อม TodoService
// output: controller ที่แปลง HTTP request เป็น service call
type todoController struct {
	utilControllers.GenericController
	todoService servicePort.TodoService
}

// NewTodoController สร้าง todo controller
// input: todoService สำหรับทำ business logic ของ todo
// output: ITodoController ที่พร้อมผูก route
func NewTodoController(todoService servicePort.TodoService) ITodoController {
	return &todoController{
		GenericController: utilControllers.NewGenericController(),
		todoService:       todoService,
	}
}
