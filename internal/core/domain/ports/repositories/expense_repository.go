package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// ExpenseRepository คือ contract สำหรับอ่าน/เขียนรายจ่ายและรวมยอดรายจ่าย
// input: domain entity, userID, id หรือช่วงเวลา spent_at
// output: domain entity/list/total หรือ error จาก adapter ที่ implement จริง
type ExpenseRepository interface {
	Create(expense *entities.Expense) (*entities.Expense, error)
	GetByID(id uint) (*entities.Expense, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Expense, error)
	ListBySpentAtBetween(userID uint, start, end time.Time) ([]*entities.Expense, error)
	SumBySpentAtBetween(userID uint, start, end time.Time) (float64, error)
	Update(expense *entities.Expense) (*entities.Expense, error)
	Delete(id uint) error
}
