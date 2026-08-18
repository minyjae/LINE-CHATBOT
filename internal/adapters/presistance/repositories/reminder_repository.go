package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type reminderRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ReminderRepository = (*reminderRepositoryImpl)(nil)

func NewReminderRepository(db *gorm.DB) *reminderRepositoryImpl {
	return &reminderRepositoryImpl{db: db}
}

func (r *reminderRepositoryImpl) Create(reminder *entities.Reminder) (*entities.Reminder, error) {
	m := models.ReminderFromEntity(reminder)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *reminderRepositoryImpl) GetByID(id uint) (*entities.Reminder, error) {
	var m models.Reminder
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *reminderRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Reminder, error) {
	var rows []models.Reminder
	if err := r.db.Where("user_id = ?", userID).
		Order("remind_at asc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return remindersToEntities(rows), nil
}

func (r *reminderRepositoryImpl) ListPendingDue(now time.Time, limit int) ([]*entities.Reminder, error) {
	var rows []models.Reminder
	if err := r.db.Where("status = ? AND remind_at <= ?", "pending", now).
		Order("remind_at asc").
		Limit(normalizeLimit(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return remindersToEntities(rows), nil
}

func (r *reminderRepositoryImpl) Update(reminder *entities.Reminder) (*entities.Reminder, error) {
	m := models.ReminderFromEntity(reminder)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *reminderRepositoryImpl) MarkSent(id uint, sentAt time.Time) error {
	return r.db.Model(&models.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":  "sent",
			"sent_at": sentAt,
		}).Error
}

func (r *reminderRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Reminder{}, id).Error
}

func remindersToEntities(rows []models.Reminder) []*entities.Reminder {
	result := make([]*entities.Reminder, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
