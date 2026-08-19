package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type CalendarEventService interface {
	List(userID uint, limit, offset int) ([]*entities.CalendarEvent, error)
	ListByDate(userID uint, date time.Time, loc *time.Location) ([]*entities.CalendarEvent, error)
}
