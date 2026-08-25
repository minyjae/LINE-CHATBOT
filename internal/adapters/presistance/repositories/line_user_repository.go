package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// lineUserRepositoryImpl เป็น GORM implementation ของ LineUserRepository
// input: สร้างจาก NewLineUserRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table line_users และคืน domain entity
type lineUserRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.LineUserRepository = (*lineUserRepositoryImpl)(nil)

// NewLineUserRepository สร้าง LINE user repository
// input: db connection ของ GORM
// output: *lineUserRepositoryImpl ที่พร้อมใช้ query LINE user
func NewLineUserRepository(db *gorm.DB) *lineUserRepositoryImpl {
	return &lineUserRepositoryImpl{db: db}
}

// Create บันทึก LINE user ใหม่ลง database
// input: lineUser domain entity ที่ต้องการบันทึก
// output: *LineUser ที่บันทึกแล้ว หรือ error จาก GORM
func (r *lineUserRepositoryImpl) Create(lineUser *entities.LineUser) (*entities.LineUser, error) {
	m := models.LineUserFromEntity(lineUser)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง LINE user จาก primary key
// input: id ของ LINE user record
// output: *LineUser ที่พบ หรือ error เช่น record not found
func (r *lineUserRepositoryImpl) GetByID(id uint) (*entities.LineUser, error) {
	var m models.LineUser
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByLineUserID ดึง LINE user จาก user id ที่ LINE ส่งมา
// input: lineUserID จาก LINE source user id
// output: *LineUser ที่พบ หรือ error เช่น record not found
func (r *lineUserRepositoryImpl) GetByLineUserID(lineUserID string) (*entities.LineUser, error) {
	var m models.LineUser
	if err := r.db.Where("line_user_id = ?", lineUserID).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง LINE account ทั้งหมดที่ผูกกับ user ในระบบ
// input: userID เจ้าของข้อมูล
// output: []*LineUser เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *lineUserRepositoryImpl) ListByUserID(userID uint) ([]*entities.LineUser, error) {
	var rows []models.LineUser
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.LineUser, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

// Update บันทึกการแก้ไข LINE user
// input: lineUser domain entity ที่มี ID และค่าใหม่
// output: *LineUser หลัง save หรือ error จาก GORM
func (r *lineUserRepositoryImpl) Update(lineUser *entities.LineUser) (*entities.LineUser, error) {
	m := models.LineUserFromEntity(lineUser)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ LINE user record ตาม id
// input: id ของ LINE user record
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *lineUserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.LineUser{}, id).Error
}
