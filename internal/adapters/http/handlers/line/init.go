package line

import (
	lineAdapter "minyjae/go-starter/internal/adapters/line"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	utilControllers "minyjae/go-starter/utils/controller"

	"github.com/gofiber/fiber/v2"
)

type ILineController interface {
	Webhook(c *fiber.Ctx) error
}

type lineController struct {
	utilControllers.GenericController
	lineWebhookService servicePort.LineWebhookService
	lineMessenger      lineAdapter.Messenger
	channelSecret      string
}

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
