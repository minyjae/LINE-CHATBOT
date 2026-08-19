package todo

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type ITodoController interface {
	List(c *fiber.Ctx) error
	ListPending(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type todoController struct {
	utilControllers.GenericController
	todoService servicePort.TodoService
}

func NewTodoController(todoService servicePort.TodoService) ITodoController {
	return &todoController{
		GenericController: utilControllers.NewGenericController(),
		todoService:       todoService,
	}
}
