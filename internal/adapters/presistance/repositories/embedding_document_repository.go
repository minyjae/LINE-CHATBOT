package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type embeddingDocumentRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.EmbeddingDocumentRepository = (*embeddingDocumentRepositoryImpl)(nil)

func NewEmbeddingDocumentRepository(db *gorm.DB) *embeddingDocumentRepositoryImpl {
	return &embeddingDocumentRepositoryImpl{db: db}
}

func (r *embeddingDocumentRepositoryImpl) Create(document *entities.EmbeddingDocument) (*entities.EmbeddingDocument, error) {
	m := models.EmbeddingDocumentFromEntity(document)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *embeddingDocumentRepositoryImpl) GetByID(id uint) (*entities.EmbeddingDocument, error) {
	var m models.EmbeddingDocument
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *embeddingDocumentRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.EmbeddingDocument, error) {
	var rows []models.EmbeddingDocument
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return embeddingDocumentsToEntities(rows), nil
}

func (r *embeddingDocumentRepositoryImpl) ListBySource(userID uint, sourceType string, sourceID *uint) ([]*entities.EmbeddingDocument, error) {
	query := r.db.Where("user_id = ? AND source_type = ?", userID, sourceType)
	if sourceID != nil {
		query = query.Where("source_id = ?", *sourceID)
	}

	var rows []models.EmbeddingDocument
	if err := query.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return embeddingDocumentsToEntities(rows), nil
}

func (r *embeddingDocumentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.EmbeddingDocument{}, id).Error
}

func embeddingDocumentsToEntities(rows []models.EmbeddingDocument) []*entities.EmbeddingDocument {
	result := make([]*entities.EmbeddingDocument, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
