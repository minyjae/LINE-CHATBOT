package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type LineUserRepository interface {
	Create(lineUser *entities.LineUser) (*entities.LineUser, error)
	GetByID(id uint) (*entities.LineUser, error)
	GetByLineUserID(lineUserID string) (*entities.LineUser, error)
	ListByUserID(userID uint) ([]*entities.LineUser, error)
	Update(lineUser *entities.LineUser) (*entities.LineUser, error)
	Delete(id uint) error
}
