package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// MessageLogRepository คือ contract สำหรับเก็บและอ่าน message log
// input: message log entity, id หรือ userID พร้อม limit/offset
// output: message log entity/list หรือ error จาก adapter ที่ implement จริง
type MessageLogRepository interface {
	Create(messageLog *entities.MessageLog) (*entities.MessageLog, error)
	GetByID(id uint) (*entities.MessageLog, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.MessageLog, error)
}
