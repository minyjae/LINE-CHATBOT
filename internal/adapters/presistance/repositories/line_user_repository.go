package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type lineUserRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.LineUserRepository = (*lineUserRepositoryImpl)(nil)

func NewLineUserRepository(db *gorm.DB) *lineUserRepositoryImpl {
	return &lineUserRepositoryImpl{db: db}
}

func (r *lineUserRepositoryImpl) Create(lineUser *entities.LineUser) (*entities.LineUser, error) {
	m := models.LineUserFromEntity(lineUser)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *lineUserRepositoryImpl) GetByID(id uint) (*entities.LineUser, error) {
	var m models.LineUser
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *lineUserRepositoryImpl) GetByLineUserID(lineUserID string) (*entities.LineUser, error) {
	var m models.LineUser
	if err := r.db.Where("line_user_id = ?", lineUserID).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *lineUserRepositoryImpl) ListByUserID(userID uint) ([]*entities.LineUser, error) {
	var rows []models.LineUser
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.LineUser, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

func (r *lineUserRepositoryImpl) Update(lineUser *entities.LineUser) (*entities.LineUser, error) {
	m := models.LineUserFromEntity(lineUser)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *lineUserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.LineUser{}, id).Error
}
