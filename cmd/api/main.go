package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"minyjae/go-starter/internal/adapters/http/routes"
	"minyjae/go-starter/internal/adapters/jobs"
	lineAdapter "minyjae/go-starter/internal/adapters/line"
	"minyjae/go-starter/internal/adapters/presistance/repositories"
	"minyjae/go-starter/internal/config"
	coreServices "minyjae/go-starter/internal/core/services"
	"minyjae/go-starter/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db := config.SetupDatabase(cfg)

	userRepo := repositories.NewUserRepository(db)
	lineUserRepo := repositories.NewLineUserRepository(db)
	messageLogRepo := repositories.NewMessageLogRepository(db)
	assistantIntentRepo := repositories.NewAssistantIntentRepository(db)
	todoRepo := repositories.NewTodoRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	calendarEventRepo := repositories.NewCalendarEventRepository(db)
	reminderRepo := repositories.NewReminderRepository(db)
	noteRepo := repositories.NewNoteRepository(db)

	assistantService := coreServices.NewAssistantServiceImpl(
		assistantIntentRepo,
		todoRepo,
		expenseRepo,
		calendarEventRepo,
		reminderRepo,
		noteRepo,
	)
	lineWebhookService := coreServices.NewLineWebhookServiceImpl(
		userRepo,
		lineUserRepo,
		messageLogRepo,
		assistantService,
	)
	todoService := coreServices.NewTodoServiceImpl(todoRepo)
	expenseService := coreServices.NewExpenseServiceImpl(expenseRepo)
	calendarEventService := coreServices.NewCalendarEventServiceImpl(calendarEventRepo)
	reminderService := coreServices.NewReminderServiceImpl(reminderRepo)
	noteService := coreServices.NewNoteServiceImpl(noteRepo)
	lineMessenger := lineAdapter.NewClient(cfg.LineChannelAccessToken)

	if cfg.ReminderWorkerEnabled {
		reminderWorker := jobs.NewReminderWorker(
			reminderRepo,
			lineUserRepo,
			messageLogRepo,
			lineMessenger,
			time.Duration(cfg.ReminderWorkerIntervalSeconds)*time.Second,
			cfg.ReminderWorkerBatchSize,
		)
		reminderWorker.Start(context.Background())
	}

	resp := utils.NewResponse()
	app := fiber.New(fiber.Config{
		ErrorHandler: resp.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsAllows,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	routes.SetupRoute(
		app,
		lineWebhookService,
		lineMessenger,
		cfg.LineChannelSecret,
		routes.DashboardServices{
			Todo:          todoService,
			Expense:       expenseService,
			CalendarEvent: calendarEventService,
			Reminder:      reminderService,
			Note:          noteService,
		},
	)
	log.Printf("Server starting on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
