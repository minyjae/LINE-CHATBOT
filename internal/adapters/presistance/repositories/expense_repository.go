package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// expenseRepositoryImpl เป็น GORM implementation ของ ExpenseRepository
// input: สร้างจาก NewExpenseRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table expenses และคืน domain entity
type expenseRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ExpenseRepository = (*expenseRepositoryImpl)(nil)

// NewExpenseRepository สร้าง expense repository
// input: db connection ของ GORM
// output: *expenseRepositoryImpl ที่พร้อมใช้ query รายจ่าย
func NewExpenseRepository(db *gorm.DB) *expenseRepositoryImpl {
	return &expenseRepositoryImpl{db: db}
}

// Create บันทึกรายจ่ายใหม่ลง database
// input: expense domain entity ที่ต้องการบันทึก
// output: *Expense ที่บันทึกแล้ว หรือ error จาก GORM
func (r *expenseRepositoryImpl) Create(expense *entities.Expense) (*entities.Expense, error) {
	m := models.ExpenseFromEntity(expense)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึงรายจ่ายจาก primary key
// input: id ของรายจ่าย
// output: *Expense ที่พบ หรือ error เช่น record not found
func (r *expenseRepositoryImpl) GetByID(id uint) (*entities.Expense, error) {
	var m models.Expense
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึงรายจ่ายของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Expense เรียงตาม spent_at/created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *expenseRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Expense, error) {
	var rows []models.Expense
	if err := r.db.Where("user_id = ?", userID).
		Order("spent_at desc, created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return expensesToEntities(rows), nil
}

// ListBySpentAtBetween ดึงรายจ่ายที่ spent_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: []*Expense ที่อยู่ในช่วงและเรียงล่าสุดก่อน หรือ error จาก GORM
func (r *expenseRepositoryImpl) ListBySpentAtBetween(userID uint, start, end time.Time) ([]*entities.Expense, error) {
	var rows []models.Expense
	if err := r.db.Where("user_id = ? AND spent_at >= ? AND spent_at < ?", userID, start, end).
		Order("spent_at desc, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return expensesToEntities(rows), nil
}

// SumBySpentAtBetween รวมยอดรายจ่ายที่ spent_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: float64 ยอดรวม amount หรือ error จาก GORM
func (r *expenseRepositoryImpl) SumBySpentAtBetween(userID uint, start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND spent_at >= ? AND spent_at < ?", userID, start, end).
		Scan(&total).Error
	return total, err
}

// Update บันทึกการแก้ไขรายจ่าย
// input: expense domain entity ที่มี ID และค่าใหม่
// output: *Expense หลัง save หรือ error จาก GORM
func (r *expenseRepositoryImpl) Update(expense *entities.Expense) (*entities.Expense, error) {
	m := models.ExpenseFromEntity(expense)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบรายจ่ายตาม id
// input: id ของรายจ่าย
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *expenseRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}

// expensesToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.Expense จาก GORM
// output: []*entities.Expense สำหรับส่งออกจาก repository
func expensesToEntities(rows []models.Expense) []*entities.Expense {
	result := make([]*entities.Expense, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
