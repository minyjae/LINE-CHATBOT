package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type AssistantIntentRepository interface {
	Create(intent *entities.AssistantIntent) (*entities.AssistantIntent, error)
	GetByID(id uint) (*entities.AssistantIntent, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.AssistantIntent, error)
	Update(intent *entities.AssistantIntent) (*entities.AssistantIntent, error)
	UpdateStatus(id uint, status string, errorMessage *string) error
}
