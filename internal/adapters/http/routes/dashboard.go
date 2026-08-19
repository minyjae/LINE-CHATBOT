package routes

import (
	calendarHandler "minyjae/go-starter/internal/adapters/http/handlers/calendar"
	expenseHandler "minyjae/go-starter/internal/adapters/http/handlers/expense"
	"minyjae/go-starter/internal/adapters/http/handlers/note"
	reminderHandler "minyjae/go-starter/internal/adapters/http/handlers/reminder"
	todoHandler "minyjae/go-starter/internal/adapters/http/handlers/todo"
	"minyjae/go-starter/internal/adapters/http/middleware"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardServices struct {
	Todo          servicePort.TodoService
	Expense       servicePort.ExpenseService
	CalendarEvent servicePort.CalendarEventService
	Reminder      servicePort.ReminderService
	Note          servicePort.NoteService
}

func DashboardRoute(app *fiber.App, services DashboardServices) {
	group := app.Group("/dashboard", middleware.JWTProtected())

	todoController := todoHandler.NewTodoController(services.Todo)
	group.Get("/todos", todoController.List)
	group.Get("/todos/pending", todoController.ListPending)

	expenseController := expenseHandler.NewExpenseController(services.Expense)
	group.Get("/expenses", expenseController.List)
	group.Get("/expenses/summary", expenseController.Summary)

	calendarController := calendarHandler.NewCalendarController(services.CalendarEvent)
	group.Get("/calendar/events", calendarController.List)
	group.Get("/calendar/events/by-date", calendarController.ListByDate)

	reminderController := reminderHandler.NewReminderController(services.Reminder)
	group.Get("/reminders", reminderController.List)

	noteController := note.NewNoteController(services.Note)
	group.Get("/notes", noteController.List)
	group.Get("/notes/search", noteController.Search)
}
