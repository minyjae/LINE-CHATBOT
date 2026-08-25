package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// IncomeSummary คือผลรวมรายรับตามช่วงเวลา
// input: สร้างจาก IncomeService.SummaryByPeriod/SummaryByMonth
// output: struct ที่มี total, currency และช่วงเวลา start/end
type IncomeSummary struct {
	Total    float64   `json:"total"`
	Currency string    `json:"currency"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

// IncomeService คือ contract business logic สำหรับรายรับ
// input: userID, income entity, id หรือช่วงเวลา
// output: income entity/list/summary หรือ error
type IncomeService interface {
	Create(userID uint, income *entities.Income) (*entities.Income, error)
	List(userID uint, limit, offset int) ([]*entities.Income, error)
	SummaryByPeriod(userID uint, start, end time.Time) (*IncomeSummary, error)
	SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*IncomeSummary, error)
	Update(userID, id uint, income *entities.Income) (*entities.Income, error)
	Delete(userID, id uint) error
}
