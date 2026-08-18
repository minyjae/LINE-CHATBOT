package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type CalendarEventRepository interface {
	Create(event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	GetByID(id uint) (*entities.CalendarEvent, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.CalendarEvent, error)
	ListByStartBetween(userID uint, start, end time.Time) ([]*entities.CalendarEvent, error)
	Update(event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	Delete(id uint) error
}
