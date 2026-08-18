package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type expenseRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ExpenseRepository = (*expenseRepositoryImpl)(nil)

func NewExpenseRepository(db *gorm.DB) *expenseRepositoryImpl {
	return &expenseRepositoryImpl{db: db}
}

func (r *expenseRepositoryImpl) Create(expense *entities.Expense) (*entities.Expense, error) {
	m := models.ExpenseFromEntity(expense)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *expenseRepositoryImpl) GetByID(id uint) (*entities.Expense, error) {
	var m models.Expense
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *expenseRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Expense, error) {
	var rows []models.Expense
	if err := r.db.Where("user_id = ?", userID).
		Order("spent_at desc, created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return expensesToEntities(rows), nil
}

func (r *expenseRepositoryImpl) ListBySpentAtBetween(userID uint, start, end time.Time) ([]*entities.Expense, error) {
	var rows []models.Expense
	if err := r.db.Where("user_id = ? AND spent_at >= ? AND spent_at < ?", userID, start, end).
		Order("spent_at desc, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return expensesToEntities(rows), nil
}

func (r *expenseRepositoryImpl) SumBySpentAtBetween(userID uint, start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND spent_at >= ? AND spent_at < ?", userID, start, end).
		Scan(&total).Error
	return total, err
}

func (r *expenseRepositoryImpl) Update(expense *entities.Expense) (*entities.Expense, error) {
	m := models.ExpenseFromEntity(expense)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *expenseRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}

func expensesToEntities(rows []models.Expense) []*entities.Expense {
	result := make([]*entities.Expense, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
