package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// CalendarEventService คือ contract business logic สำหรับ calendar event
// input: userID, calendar event entity, id, date หรือ timezone
// output: calendar event entity/list หรือ error
type CalendarEventService interface {
	Create(userID uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	List(userID uint, limit, offset int) ([]*entities.CalendarEvent, error)
	ListByDate(userID uint, date time.Time, loc *time.Location) ([]*entities.CalendarEvent, error)
	Update(userID, id uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	Delete(userID, id uint) error
}
