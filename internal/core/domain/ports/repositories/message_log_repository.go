package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type MessageLogRepository interface {
	Create(messageLog *entities.MessageLog) (*entities.MessageLog, error)
	GetByID(id uint) (*entities.MessageLog, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.MessageLog, error)
}
