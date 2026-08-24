package services

import "minyjae/go-starter/internal/core/domain/entities"

type ReminderService interface {
	Create(userID uint, reminder *entities.Reminder) (*entities.Reminder, error)
	List(userID uint, limit, offset int) ([]*entities.Reminder, error)
	Update(userID, id uint, reminder *entities.Reminder) (*entities.Reminder, error)
	Delete(userID, id uint) error
}
