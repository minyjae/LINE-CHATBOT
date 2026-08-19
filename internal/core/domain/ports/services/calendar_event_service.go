package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type CalendarEventService interface {
	Create(userID uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	List(userID uint, limit, offset int) ([]*entities.CalendarEvent, error)
	ListByDate(userID uint, date time.Time, loc *time.Location) ([]*entities.CalendarEvent, error)
	Update(userID, id uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	Delete(userID, id uint) error
}
