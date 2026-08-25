package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// incomeRepositoryImpl เป็น GORM implementation ของ IncomeRepository
// input: สร้างจาก NewIncomeRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table incomes และคืน domain entity
type incomeRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.IncomeRepository = (*incomeRepositoryImpl)(nil)

// NewIncomeRepository สร้าง income repository
// input: db connection ของ GORM
// output: *incomeRepositoryImpl ที่พร้อมใช้ query รายรับ
func NewIncomeRepository(db *gorm.DB) *incomeRepositoryImpl {
	return &incomeRepositoryImpl{db: db}
}

// Create บันทึกรายรับใหม่ลง database
// input: income domain entity ที่ต้องการบันทึก
// output: *Income ที่บันทึกแล้ว หรือ error จาก GORM
func (r *incomeRepositoryImpl) Create(income *entities.Income) (*entities.Income, error) {
	m := models.IncomeFromEntity(income)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึงรายรับจาก primary key
// input: id ของรายรับ
// output: *Income ที่พบ หรือ error เช่น record not found
func (r *incomeRepositoryImpl) GetByID(id uint) (*entities.Income, error) {
	var m models.Income
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึงรายรับของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Income เรียงตาม received_at/created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *incomeRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Income, error) {
	var rows []models.Income
	if err := r.db.Where("user_id = ?", userID).
		Order("received_at desc, created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return incomesToEntities(rows), nil
}

// ListByReceivedAtBetween ดึงรายรับที่ received_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: []*Income ที่อยู่ในช่วงและเรียงล่าสุดก่อน หรือ error จาก GORM
func (r *incomeRepositoryImpl) ListByReceivedAtBetween(userID uint, start, end time.Time) ([]*entities.Income, error) {
	var rows []models.Income
	if err := r.db.Where("user_id = ? AND received_at >= ? AND received_at < ?", userID, start, end).
		Order("received_at desc, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return incomesToEntities(rows), nil
}

// SumByReceivedAtBetween รวมยอดรายรับที่ received_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: float64 ยอดรวม amount หรือ error จาก GORM
func (r *incomeRepositoryImpl) SumByReceivedAtBetween(userID uint, start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Income{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND received_at >= ? AND received_at < ?", userID, start, end).
		Scan(&total).Error
	return total, err
}

// Update บันทึกการแก้ไขรายรับ
// input: income domain entity ที่มี ID และค่าใหม่
// output: *Income หลัง save หรือ error จาก GORM
func (r *incomeRepositoryImpl) Update(income *entities.Income) (*entities.Income, error) {
	m := models.IncomeFromEntity(income)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบรายรับตาม id
// input: id ของรายรับ
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *incomeRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Income{}, id).Error
}

// incomesToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.Income จาก GORM
// output: []*entities.Income สำหรับส่งออกจาก repository
func incomesToEntities(rows []models.Income) []*entities.Income {
	result := make([]*entities.Income, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
