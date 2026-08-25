package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// CalendarEventRepository คือ contract สำหรับอ่าน/เขียน calendar event
// input: domain entity, userID, id หรือช่วงเวลา start/end
// output: domain entity/list หรือ error จาก adapter ที่ implement จริง
type CalendarEventRepository interface {
	Create(event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	GetByID(id uint) (*entities.CalendarEvent, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.CalendarEvent, error)
	ListByStartBetween(userID uint, start, end time.Time) ([]*entities.CalendarEvent, error)
	Update(event *entities.CalendarEvent) (*entities.CalendarEvent, error)
	Delete(id uint) error
}
