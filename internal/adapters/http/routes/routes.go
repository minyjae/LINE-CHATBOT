package routes

import (
	lineAdapter "minyjae/go-starter/internal/adapters/line"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoute(
	app *fiber.App,
	lineWebhookService servicePort.LineWebhookService,
	lineMessenger lineAdapter.Messenger,
	lineChannelSecret string,
	dashboardServices DashboardServices,
) {
	LineRoute(app, lineWebhookService, lineMessenger, lineChannelSecret)
	DashboardRoute(app, dashboardServices)
}
