package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// EmbeddingDocumentRepository คือ contract สำหรับอ่าน/เขียนเอกสารที่ใช้ทำ embedding
// input: embedding document entity, id, userID หรือ source filter
// output: embedding document entity/list หรือ error จาก adapter ที่ implement จริง
type EmbeddingDocumentRepository interface {
	Create(document *entities.EmbeddingDocument) (*entities.EmbeddingDocument, error)
	GetByID(id uint) (*entities.EmbeddingDocument, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.EmbeddingDocument, error)
	ListBySource(userID uint, sourceType string, sourceID *uint) ([]*entities.EmbeddingDocument, error)
	Delete(id uint) error
}
