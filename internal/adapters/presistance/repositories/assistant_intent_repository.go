package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// assistantIntentRepositoryImpl เป็น GORM implementation ของ AssistantIntentRepository
// input: สร้างจาก NewAssistantIntentRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table assistant_intents และคืน domain entity
type assistantIntentRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.AssistantIntentRepository = (*assistantIntentRepositoryImpl)(nil)

// NewAssistantIntentRepository สร้าง assistant intent repository
// input: db connection ของ GORM
// output: *assistantIntentRepositoryImpl ที่พร้อมใช้ query assistant intent
func NewAssistantIntentRepository(db *gorm.DB) *assistantIntentRepositoryImpl {
	return &assistantIntentRepositoryImpl{db: db}
}

// Create บันทึก assistant intent ใหม่ลง database
// input: intent domain entity ที่ต้องการบันทึก
// output: *AssistantIntent ที่บันทึกแล้ว หรือ error จาก GORM
func (r *assistantIntentRepositoryImpl) Create(intent *entities.AssistantIntent) (*entities.AssistantIntent, error) {
	m := models.AssistantIntentFromEntity(intent)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง assistant intent จาก primary key
// input: id ของ assistant intent
// output: *AssistantIntent ที่พบ หรือ error เช่น record not found
func (r *assistantIntentRepositoryImpl) GetByID(id uint) (*entities.AssistantIntent, error) {
	var m models.AssistantIntent
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง assistant intent history ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*AssistantIntent เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *assistantIntentRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.AssistantIntent, error) {
	var rows []models.AssistantIntent
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.AssistantIntent, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

// Update บันทึกการแก้ไข assistant intent
// input: intent domain entity ที่มี ID และค่าใหม่
// output: *AssistantIntent หลัง save หรือ error จาก GORM
func (r *assistantIntentRepositoryImpl) Update(intent *entities.AssistantIntent) (*entities.AssistantIntent, error) {
	m := models.AssistantIntentFromEntity(intent)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// UpdateStatus แก้ไขเฉพาะ status และ error message ของ assistant intent
// input: id ของ intent, status ใหม่, errorMessage pointer ที่อาจเป็น nil
// output: error nil เมื่อ update สำเร็จ หรือ error จาก GORM
func (r *assistantIntentRepositoryImpl) UpdateStatus(id uint, status string, errorMessage *string) error {
	return r.db.Model(&models.AssistantIntent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        status,
			"error_message": errorMessage,
		}).Error
}
