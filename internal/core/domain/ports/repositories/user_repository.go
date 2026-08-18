package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type UserRepository interface {
	Create(user *entities.User) (*entities.User, error)
	GetByID(id uint) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(id uint) error
}
