package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// noteRepositoryImpl เป็น GORM implementation ของ NoteRepository
// input: สร้างจาก NewNoteRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table notes และคืน domain entity
type noteRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.NoteRepository = (*noteRepositoryImpl)(nil)

// NewNoteRepository สร้าง note repository
// input: db connection ของ GORM
// output: *noteRepositoryImpl ที่พร้อมใช้ query note
func NewNoteRepository(db *gorm.DB) *noteRepositoryImpl {
	return &noteRepositoryImpl{db: db}
}

// Create บันทึก note ใหม่ลง database
// input: note domain entity ที่ต้องการบันทึก
// output: *Note ที่บันทึกแล้ว หรือ error จาก GORM
func (r *noteRepositoryImpl) Create(note *entities.Note) (*entities.Note, error) {
	m := models.NoteFromEntity(note)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง note จาก primary key
// input: id ของ note
// output: *Note ที่พบ หรือ error เช่น record not found
func (r *noteRepositoryImpl) GetByID(id uint) (*entities.Note, error) {
	var m models.Note
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง note ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Note เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *noteRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Note, error) {
	var rows []models.Note
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return notesToEntities(rows), nil
}

// SearchByContent ค้น note ด้วยข้อความใน content
// input: userID เจ้าของข้อมูล, query คำค้นหา, limit จำนวนสูงสุด
// output: []*Note ที่ content ตรงกับ query แบบ ILIKE หรือ error จาก GORM
func (r *noteRepositoryImpl) SearchByContent(userID uint, query string, limit int) ([]*entities.Note, error) {
	var rows []models.Note
	if err := r.db.Where("user_id = ? AND content ILIKE ?", userID, "%"+query+"%").
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return notesToEntities(rows), nil
}

// Update บันทึกการแก้ไข note
// input: note domain entity ที่มี ID และค่าใหม่
// output: *Note หลัง save หรือ error จาก GORM
func (r *noteRepositoryImpl) Update(note *entities.Note) (*entities.Note, error) {
	m := models.NoteFromEntity(note)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ note ตาม id
// input: id ของ note
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *noteRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}

// notesToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.Note จาก GORM
// output: []*entities.Note สำหรับส่งออกจาก repository
func notesToEntities(rows []models.Note) []*entities.Note {
	result := make([]*entities.Note, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
