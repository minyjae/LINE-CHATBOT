package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type ExpenseSummary struct {
	Total    float64   `json:"total"`
	Currency string    `json:"currency"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type ExpenseService interface {
	List(userID uint, limit, offset int) ([]*entities.Expense, error)
	SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*ExpenseSummary, error)
}
