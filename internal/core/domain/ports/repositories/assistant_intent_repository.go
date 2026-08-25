package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// AssistantIntentRepository คือ contract สำหรับเก็บ intent ที่ assistant ตีความได้
// input: intent entity, id, userID หรือ status/errorMessage
// output: intent entity/list หรือ error จาก adapter ที่ implement จริง
type AssistantIntentRepository interface {
	Create(intent *entities.AssistantIntent) (*entities.AssistantIntent, error)
	GetByID(id uint) (*entities.AssistantIntent, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.AssistantIntent, error)
	Update(intent *entities.AssistantIntent) (*entities.AssistantIntent, error)
	UpdateStatus(id uint, status string, errorMessage *string) error
}
