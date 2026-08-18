package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type integrationAccountRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.IntegrationAccountRepository = (*integrationAccountRepositoryImpl)(nil)

func NewIntegrationAccountRepository(db *gorm.DB) *integrationAccountRepositoryImpl {
	return &integrationAccountRepositoryImpl{db: db}
}

func (r *integrationAccountRepositoryImpl) Create(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error) {
	m := models.IntegrationAccountFromEntity(account)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *integrationAccountRepositoryImpl) GetByID(id uint) (*entities.IntegrationAccount, error) {
	var m models.IntegrationAccount
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *integrationAccountRepositoryImpl) GetByUserIDAndProvider(userID uint, provider string) (*entities.IntegrationAccount, error) {
	var m models.IntegrationAccount
	if err := r.db.Where("user_id = ? AND provider = ?", userID, provider).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *integrationAccountRepositoryImpl) ListByUserID(userID uint) ([]*entities.IntegrationAccount, error) {
	var rows []models.IntegrationAccount
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.IntegrationAccount, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

func (r *integrationAccountRepositoryImpl) Update(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error) {
	m := models.IntegrationAccountFromEntity(account)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *integrationAccountRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.IntegrationAccount{}, id).Error
}
