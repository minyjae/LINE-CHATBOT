package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// messageLogRepositoryImpl เป็น GORM implementation ของ MessageLogRepository
// input: สร้างจาก NewMessageLogRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table message_logs และคืน domain entity
type messageLogRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.MessageLogRepository = (*messageLogRepositoryImpl)(nil)

// NewMessageLogRepository สร้าง message log repository
// input: db connection ของ GORM
// output: *messageLogRepositoryImpl ที่พร้อมใช้ query message log
func NewMessageLogRepository(db *gorm.DB) *messageLogRepositoryImpl {
	return &messageLogRepositoryImpl{db: db}
}

// Create บันทึก message log ใหม่ลง database
// input: messageLog domain entity ที่ต้องการบันทึก
// output: *MessageLog ที่บันทึกแล้ว หรือ error จาก GORM
func (r *messageLogRepositoryImpl) Create(messageLog *entities.MessageLog) (*entities.MessageLog, error) {
	m := models.MessageLogFromEntity(messageLog)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง message log จาก primary key
// input: id ของ message log
// output: *MessageLog ที่พบ หรือ error เช่น record not found
func (r *messageLogRepositoryImpl) GetByID(id uint) (*entities.MessageLog, error) {
	var m models.MessageLog
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง message log ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*MessageLog เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *messageLogRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.MessageLog, error) {
	var rows []models.MessageLog
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.MessageLog, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}
