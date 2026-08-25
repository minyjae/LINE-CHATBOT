package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// ExpenseSummary คือผลรวมรายจ่ายตามช่วงเวลา
// input: สร้างจาก ExpenseService.SummaryByPeriod/SummaryByMonth
// output: struct ที่มี total, currency และช่วงเวลา start/end
type ExpenseSummary struct {
	Total    float64   `json:"total"`
	Currency string    `json:"currency"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

// ExpenseService คือ contract business logic สำหรับรายจ่าย
// input: userID, expense entity, id หรือช่วงเวลา
// output: expense entity/list/summary หรือ error
type ExpenseService interface {
	Create(userID uint, expense *entities.Expense) (*entities.Expense, error)
	List(userID uint, limit, offset int) ([]*entities.Expense, error)
	SummaryByPeriod(userID uint, start, end time.Time) (*ExpenseSummary, error)
	SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*ExpenseSummary, error)
	Update(userID, id uint, expense *entities.Expense) (*entities.Expense, error)
	Delete(userID, id uint) error
}
