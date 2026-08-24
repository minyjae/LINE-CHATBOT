package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type reminderService struct {
	repo repoPort.ReminderRepository
}

var _ servicePort.ReminderService = (*reminderService)(nil)

func NewReminderServiceImpl(repo repoPort.ReminderRepository) *reminderService {
	return &reminderService{repo: repo}
}

func (s *reminderService) Create(userID uint, reminder *entities.Reminder) (*entities.Reminder, error) {
	now := time.Now()
	reminder.UserID = userID
	reminder.Status = defaultString(reminder.Status, "pending")
	reminder.CreatedAt = now
	reminder.UpdatedAt = now
	return s.repo.Create(reminder)
}

func (s *reminderService) List(userID uint, limit, offset int) ([]*entities.Reminder, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *reminderService) Update(userID, id uint, reminder *entities.Reminder) (*entities.Reminder, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	reminder.ID = id
	reminder.UserID = userID
	reminder.CreatedAt = current.CreatedAt
	reminder.UpdatedAt = time.Now()
	if reminder.Status == "" {
		reminder.Status = current.Status
	}
	return s.repo.Update(reminder)
}

func (s *reminderService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
