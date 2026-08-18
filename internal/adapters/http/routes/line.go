package routes

import (
	lineHandler "minyjae/go-starter/internal/adapters/http/handlers/line"
	lineAdapter "minyjae/go-starter/internal/adapters/line"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"

	"github.com/gofiber/fiber/v2"
)

func LineRoute(
	app *fiber.App,
	lineWebhookService servicePort.LineWebhookService,
	lineMessenger lineAdapter.Messenger,
	channelSecret string,
) {
	handler := lineHandler.NewLineController(lineWebhookService, lineMessenger, channelSecret)
	app.Post("/webhooks/line", handler.Webhook)
}
