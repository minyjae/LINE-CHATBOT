package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type calendarEventService struct {
	repo repoPort.CalendarEventRepository
}

var _ servicePort.CalendarEventService = (*calendarEventService)(nil)

func NewCalendarEventServiceImpl(repo repoPort.CalendarEventRepository) *calendarEventService {
	return &calendarEventService{repo: repo}
}

func (s *calendarEventService) List(userID uint, limit, offset int) ([]*entities.CalendarEvent, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *calendarEventService) ListByDate(userID uint, date time.Time, loc *time.Location) ([]*entities.CalendarEvent, error) {
	if loc == nil {
		loc = time.Local
	}
	localDate := date.In(loc)
	start := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)
	return s.repo.ListByStartBetween(userID, start, end)
}
