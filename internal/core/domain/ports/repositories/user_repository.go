package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// UserRepository คือ contract สำหรับอ่าน/เขียน user
// input: user entity, id หรือ email
// output: user entity หรือ error จาก adapter ที่ implement จริง
type UserRepository interface {
	Create(user *entities.User) (*entities.User, error)
	GetByID(id uint) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(id uint) error
}
