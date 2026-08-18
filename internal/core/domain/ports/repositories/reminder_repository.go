package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type ReminderRepository interface {
	Create(reminder *entities.Reminder) (*entities.Reminder, error)
	GetByID(id uint) (*entities.Reminder, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Reminder, error)
	ListPendingDue(now time.Time, limit int) ([]*entities.Reminder, error)
	Update(reminder *entities.Reminder) (*entities.Reminder, error)
	MarkSent(id uint, sentAt time.Time) error
	Delete(id uint) error
}
