package line

import (
	lineAdapter "minyjae/go-starter/internal/adapters/line"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

// ILineController กำหนด HTTP handler ของ LINE webhook
// input: Fiber context จาก route webhook
// output: error จาก Fiber handler โดย response ถูกเขียนผ่าน GenericController
type ILineController interface {
	Webhook(c *fiber.Ctx) error
}

// lineController เก็บ dependency สำหรับรับ webhook และส่งข้อความตอบกลับ LINE
// input: สร้างจาก NewLineController พร้อม LineWebhookService, Messenger และ channelSecret
// output: controller ที่ verify payload, เรียก service, และ reply กลับ LINE
type lineController struct {
	utilControllers.GenericController
	lineWebhookService servicePort.LineWebhookService
	lineMessenger      lineAdapter.Messenger
	channelSecret      string
}

// NewLineController สร้าง LINE controller
// input: lineWebhookService สำหรับ business logic, lineMessenger สำหรับ reply, channelSecret สำหรับ verify signature
// output: ILineController ที่พร้อมผูก route webhook
func NewLineController(
	lineWebhookService servicePort.LineWebhookService,
	lineMessenger lineAdapter.Messenger,
	channelSecret string,
) ILineController {
	return &lineController{
		GenericController:  utilControllers.NewGenericController(),
		lineWebhookService: lineWebhookService,
		lineMessenger:      lineMessenger,
		channelSecret:      channelSecret,
	}
}
