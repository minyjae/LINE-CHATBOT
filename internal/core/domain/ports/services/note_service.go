package services

import "minyjae/go-starter/internal/core/domain/entities"

type NoteService interface {
	List(userID uint, limit, offset int) ([]*entities.Note, error)
	Search(userID uint, query string, limit int) ([]*entities.Note, error)
}
