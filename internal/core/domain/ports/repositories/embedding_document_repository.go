package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type EmbeddingDocumentRepository interface {
	Create(document *entities.EmbeddingDocument) (*entities.EmbeddingDocument, error)
	GetByID(id uint) (*entities.EmbeddingDocument, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.EmbeddingDocument, error)
	ListBySource(userID uint, sourceType string, sourceID *uint) ([]*entities.EmbeddingDocument, error)
	Delete(id uint) error
}
