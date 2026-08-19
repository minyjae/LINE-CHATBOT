package services

import "minyjae/go-starter/internal/core/domain/entities"

type ReminderService interface {
	List(userID uint, limit, offset int) ([]*entities.Reminder, error)
}
