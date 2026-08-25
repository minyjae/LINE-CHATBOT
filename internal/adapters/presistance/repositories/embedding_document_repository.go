package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// embeddingDocumentRepositoryImpl เป็น GORM implementation ของ EmbeddingDocumentRepository
// input: สร้างจาก NewEmbeddingDocumentRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table embedding_documents และคืน domain entity
type embeddingDocumentRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.EmbeddingDocumentRepository = (*embeddingDocumentRepositoryImpl)(nil)

// NewEmbeddingDocumentRepository สร้าง embedding document repository
// input: db connection ของ GORM
// output: *embeddingDocumentRepositoryImpl ที่พร้อมใช้ query embedding document
func NewEmbeddingDocumentRepository(db *gorm.DB) *embeddingDocumentRepositoryImpl {
	return &embeddingDocumentRepositoryImpl{db: db}
}

// Create บันทึก embedding document ใหม่ลง database
// input: document domain entity ที่ต้องการบันทึก
// output: *EmbeddingDocument ที่บันทึกแล้ว หรือ error จาก GORM
func (r *embeddingDocumentRepositoryImpl) Create(document *entities.EmbeddingDocument) (*entities.EmbeddingDocument, error) {
	m := models.EmbeddingDocumentFromEntity(document)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง embedding document จาก primary key
// input: id ของ embedding document
// output: *EmbeddingDocument ที่พบ หรือ error เช่น record not found
func (r *embeddingDocumentRepositoryImpl) GetByID(id uint) (*entities.EmbeddingDocument, error) {
	var m models.EmbeddingDocument
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง embedding document ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*EmbeddingDocument เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
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

// ListBySource ดึง embedding document ตาม source type และ optional source id
// input: userID เจ้าของข้อมูล, sourceType ประเภทแหล่งที่มา, sourceID id ของแหล่งที่มาเมื่อมี
// output: []*EmbeddingDocument ที่ตรง source หรือ error จาก GORM
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

// Delete ลบ embedding document ตาม id
// input: id ของ embedding document
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *embeddingDocumentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.EmbeddingDocument{}, id).Error
}

// embeddingDocumentsToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.EmbeddingDocument จาก GORM
// output: []*entities.EmbeddingDocument สำหรับส่งออกจาก repository
func embeddingDocumentsToEntities(rows []models.EmbeddingDocument) []*entities.EmbeddingDocument {
	result := make([]*entities.EmbeddingDocument, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
