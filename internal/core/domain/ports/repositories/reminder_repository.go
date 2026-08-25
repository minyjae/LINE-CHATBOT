package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

// ReminderRepository คือ contract สำหรับอ่าน/เขียน reminder และสถานะการส่งเตือน
// input: domain entity, userID, id, now/upcomingUntil หรือ sentAt
// output: domain entity/list หรือ error จาก adapter ที่ implement จริง
type ReminderRepository interface {
	Create(reminder *entities.Reminder) (*entities.Reminder, error)
	GetByID(id uint) (*entities.Reminder, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Reminder, error)
	ListPendingDue(now time.Time, limit int) ([]*entities.Reminder, error)
	ListPendingDueOrUpcoming(now, upcomingUntil time.Time, limit int) ([]*entities.Reminder, error)
	Update(reminder *entities.Reminder) (*entities.Reminder, error)
	MarkPreSent(id uint, sentAt time.Time) error
	MarkSent(id uint, sentAt time.Time) error
	Delete(id uint) error
}
