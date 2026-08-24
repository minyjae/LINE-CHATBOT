package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type incomeRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.IncomeRepository = (*incomeRepositoryImpl)(nil)

func NewIncomeRepository(db *gorm.DB) *incomeRepositoryImpl {
	return &incomeRepositoryImpl{db: db}
}

func (r *incomeRepositoryImpl) Create(income *entities.Income) (*entities.Income, error) {
	m := models.IncomeFromEntity(income)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *incomeRepositoryImpl) GetByID(id uint) (*entities.Income, error) {
	var m models.Income
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *incomeRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Income, error) {
	var rows []models.Income
	if err := r.db.Where("user_id = ?", userID).
		Order("received_at desc, created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return incomesToEntities(rows), nil
}

func (r *incomeRepositoryImpl) ListByReceivedAtBetween(userID uint, start, end time.Time) ([]*entities.Income, error) {
	var rows []models.Income
	if err := r.db.Where("user_id = ? AND received_at >= ? AND received_at < ?", userID, start, end).
		Order("received_at desc, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return incomesToEntities(rows), nil
}

func (r *incomeRepositoryImpl) SumByReceivedAtBetween(userID uint, start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Income{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND received_at >= ? AND received_at < ?", userID, start, end).
		Scan(&total).Error
	return total, err
}

func (r *incomeRepositoryImpl) Update(income *entities.Income) (*entities.Income, error) {
	m := models.IncomeFromEntity(income)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *incomeRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Income{}, id).Error
}

func incomesToEntities(rows []models.Income) []*entities.Income {
	result := make([]*entities.Income, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
