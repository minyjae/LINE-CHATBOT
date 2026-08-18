package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type IntegrationAccountRepository interface {
	Create(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error)
	GetByID(id uint) (*entities.IntegrationAccount, error)
	GetByUserIDAndProvider(userID uint, provider string) (*entities.IntegrationAccount, error)
	ListByUserID(userID uint) ([]*entities.IntegrationAccount, error)
	Update(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error)
	Delete(id uint) error
}
