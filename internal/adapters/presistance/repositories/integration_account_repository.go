package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// integrationAccountRepositoryImpl เป็น GORM implementation ของ IntegrationAccountRepository
// input: สร้างจาก NewIntegrationAccountRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table integration_accounts และคืน domain entity
type integrationAccountRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.IntegrationAccountRepository = (*integrationAccountRepositoryImpl)(nil)

// NewIntegrationAccountRepository สร้าง integration account repository
// input: db connection ของ GORM
// output: *integrationAccountRepositoryImpl ที่พร้อมใช้ query integration account
func NewIntegrationAccountRepository(db *gorm.DB) *integrationAccountRepositoryImpl {
	return &integrationAccountRepositoryImpl{db: db}
}

// Create บันทึก integration account ใหม่ลง database
// input: account domain entity ที่ต้องการบันทึก
// output: *IntegrationAccount ที่บันทึกแล้ว หรือ error จาก GORM
func (r *integrationAccountRepositoryImpl) Create(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error) {
	m := models.IntegrationAccountFromEntity(account)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง integration account จาก primary key
// input: id ของ integration account
// output: *IntegrationAccount ที่พบ หรือ error เช่น record not found
func (r *integrationAccountRepositoryImpl) GetByID(id uint) (*entities.IntegrationAccount, error) {
	var m models.IntegrationAccount
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByUserIDAndProvider ดึง integration account ของ user ตาม provider
// input: userID เจ้าของข้อมูล, provider ชื่อ provider เช่น google_calendar
// output: *IntegrationAccount ที่พบ หรือ error เช่น record not found
func (r *integrationAccountRepositoryImpl) GetByUserIDAndProvider(userID uint, provider string) (*entities.IntegrationAccount, error) {
	var m models.IntegrationAccount
	if err := r.db.Where("user_id = ? AND provider = ?", userID, provider).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง integration account ทั้งหมดของ user
// input: userID เจ้าของข้อมูล
// output: []*IntegrationAccount เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *integrationAccountRepositoryImpl) ListByUserID(userID uint) ([]*entities.IntegrationAccount, error) {
	var rows []models.IntegrationAccount
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.IntegrationAccount, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

// Update บันทึกการแก้ไข integration account
// input: account domain entity ที่มี ID และค่าใหม่
// output: *IntegrationAccount หลัง save หรือ error จาก GORM
func (r *integrationAccountRepositoryImpl) Update(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error) {
	m := models.IntegrationAccountFromEntity(account)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ integration account ตาม id
// input: id ของ integration account
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *integrationAccountRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.IntegrationAccount{}, id).Error
}
