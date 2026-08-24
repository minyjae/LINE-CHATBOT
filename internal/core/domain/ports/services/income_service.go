package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type IncomeSummary struct {
	Total    float64   `json:"total"`
	Currency string    `json:"currency"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type IncomeService interface {
	Create(userID uint, income *entities.Income) (*entities.Income, error)
	List(userID uint, limit, offset int) ([]*entities.Income, error)
	SummaryByPeriod(userID uint, start, end time.Time) (*IncomeSummary, error)
	SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*IncomeSummary, error)
	Update(userID, id uint, income *entities.Income) (*entities.Income, error)
	Delete(userID, id uint) error
}
