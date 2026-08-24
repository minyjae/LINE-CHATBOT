package note

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type INoteController interface {
	List(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type noteController struct {
	utilControllers.GenericController
	noteService servicePort.NoteService
}

func NewNoteController(noteService servicePort.NoteService) INoteController {
	return &noteController{
		GenericController: utilControllers.NewGenericController(),
		noteService:       noteService,
	}
}
