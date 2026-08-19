package services

import (
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

func (s *reminderService) List(userID uint, limit, offset int) ([]*entities.Reminder, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}
