package routes

import (
	calendarHandler "minyjae/go-starter/internal/adapters/http/handlers/calendar"
	expenseHandler "minyjae/go-starter/internal/adapters/http/handlers/expense"
	financeHandler "minyjae/go-starter/internal/adapters/http/handlers/finance"
	incomeHandler "minyjae/go-starter/internal/adapters/http/handlers/income"
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
	Income        servicePort.IncomeService
	CalendarEvent servicePort.CalendarEventService
	Reminder      servicePort.ReminderService
	Note          servicePort.NoteService
}

func DashboardRoute(app *fiber.App, services DashboardServices) {
	group := app.Group("/dashboard", middleware.JWTProtected())

	todoController := todoHandler.NewTodoController(services.Todo)
	group.Get("/todos", todoController.List)
	group.Get("/todos/pending", todoController.ListPending)
	group.Post("/todos", todoController.Create)
	group.Put("/todos/:id", todoController.Update)
	group.Delete("/todos/:id", todoController.Delete)

	expenseController := expenseHandler.NewExpenseController(services.Expense)
	group.Get("/expenses", expenseController.List)
	group.Get("/expenses/summary", expenseController.Summary)
	group.Post("/expenses", expenseController.Create)
	group.Put("/expenses/:id", expenseController.Update)
	group.Delete("/expenses/:id", expenseController.Delete)

	incomeController := incomeHandler.NewIncomeController(services.Income)
	group.Get("/incomes", incomeController.List)
	group.Get("/incomes/summary", incomeController.Summary)
	group.Post("/incomes", incomeController.Create)
	group.Put("/incomes/:id", incomeController.Update)
	group.Delete("/incomes/:id", incomeController.Delete)

	financeController := financeHandler.NewFinanceController(services.Expense, services.Income)
	group.Get("/finance/summary", financeController.Summary)

	calendarController := calendarHandler.NewCalendarController(services.CalendarEvent)
	group.Get("/calendar/events", calendarController.List)
	group.Get("/calendar/events/by-date", calendarController.ListByDate)
	group.Post("/calendar/events", calendarController.Create)
	group.Put("/calendar/events/:id", calendarController.Update)
	group.Delete("/calendar/events/:id", calendarController.Delete)

	reminderController := reminderHandler.NewReminderController(services.Reminder)
	group.Get("/reminders", reminderController.List)
	group.Post("/reminders", reminderController.Create)
	group.Put("/reminders/:id", reminderController.Update)
	group.Delete("/reminders/:id", reminderController.Delete)

	noteController := note.NewNoteController(services.Note)
	group.Get("/notes", noteController.List)
	group.Get("/notes/search", noteController.Search)
	group.Post("/notes", noteController.Create)
	group.Put("/notes/:id", noteController.Update)
	group.Delete("/notes/:id", noteController.Delete)
}
