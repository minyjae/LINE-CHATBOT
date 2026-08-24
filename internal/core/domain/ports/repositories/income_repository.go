package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type IncomeRepository interface {
	Create(income *entities.Income) (*entities.Income, error)
	GetByID(id uint) (*entities.Income, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Income, error)
	ListByReceivedAtBetween(userID uint, start, end time.Time) ([]*entities.Income, error)
	SumByReceivedAtBetween(userID uint, start, end time.Time) (float64, error)
	Update(income *entities.Income) (*entities.Income, error)
	Delete(id uint) error
}
