package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.UserRepository = (*userRepositoryImpl)(nil)

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *entities.User) (*entities.User, error) {
	m := models.FromEntity(user)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *userRepositoryImpl) GetByID(id uint) (*entities.User, error) {
	var m models.User
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*entities.User, error) {
	var m models.User
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *userRepositoryImpl) Update(user *entities.User) (*entities.User, error) {
	m := models.FromEntity(user)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
