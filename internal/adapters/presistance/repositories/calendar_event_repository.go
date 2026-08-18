package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type calendarEventRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.CalendarEventRepository = (*calendarEventRepositoryImpl)(nil)

func NewCalendarEventRepository(db *gorm.DB) *calendarEventRepositoryImpl {
	return &calendarEventRepositoryImpl{db: db}
}

func (r *calendarEventRepositoryImpl) Create(event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	m := models.CalendarEventFromEntity(event)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *calendarEventRepositoryImpl) GetByID(id uint) (*entities.CalendarEvent, error) {
	var m models.CalendarEvent
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *calendarEventRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.CalendarEvent, error) {
	var rows []models.CalendarEvent
	if err := r.db.Where("user_id = ?", userID).
		Order("start_at asc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return calendarEventsToEntities(rows), nil
}

func (r *calendarEventRepositoryImpl) ListByStartBetween(userID uint, start, end time.Time) ([]*entities.CalendarEvent, error) {
	var rows []models.CalendarEvent
	if err := r.db.Where("user_id = ? AND start_at >= ? AND start_at < ?", userID, start, end).
		Order("start_at asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return calendarEventsToEntities(rows), nil
}

func (r *calendarEventRepositoryImpl) Update(event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	m := models.CalendarEventFromEntity(event)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *calendarEventRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.CalendarEvent{}, id).Error
}

func calendarEventsToEntities(rows []models.CalendarEvent) []*entities.CalendarEvent {
	result := make([]*entities.CalendarEvent, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
