package note

import (
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// INoteController กำหนด HTTP handler ทั้งหมดของ note endpoint
// input: Fiber context จาก route ที่เรียกแต่ละ method
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type INoteController interface {
	List(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// noteController เก็บ dependency สำหรับ note HTTP handler
// input: สร้างจาก NewNoteController พร้อม NoteService
// output: controller ที่แปลง HTTP request เป็น service call
type noteController struct {
	utilControllers.GenericController
	noteService servicePort.NoteService
}

// NewNoteController สร้าง note controller
// input: noteService สำหรับทำ business logic ของ note
// output: INoteController ที่พร้อมผูก route
func NewNoteController(noteService servicePort.NoteService) INoteController {
	return &noteController{
		GenericController: utilControllers.NewGenericController(),
		noteService:       noteService,
	}
}
