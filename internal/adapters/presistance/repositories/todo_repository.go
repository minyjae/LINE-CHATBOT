package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// todoRepositoryImpl เป็น GORM implementation ของ TodoRepository
// input: สร้างจาก NewTodoRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table todos และคืน domain entity
type todoRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.TodoRepository = (*todoRepositoryImpl)(nil)

// NewTodoRepository สร้าง todo repository
// input: db connection ของ GORM
// output: *todoRepositoryImpl ที่พร้อมใช้ query todo
func NewTodoRepository(db *gorm.DB) *todoRepositoryImpl {
	return &todoRepositoryImpl{db: db}
}

// Create บันทึก todo ใหม่ลง database
// input: todo domain entity ที่ต้องการบันทึก
// output: *Todo ที่บันทึกแล้ว หรือ error จาก GORM
func (r *todoRepositoryImpl) Create(todo *entities.Todo) (*entities.Todo, error) {
	m := models.TodoFromEntity(todo)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง todo จาก primary key
// input: id ของ todo
// output: *Todo ที่พบ หรือ error เช่น record not found
func (r *todoRepositoryImpl) GetByID(id uint) (*entities.Todo, error) {
	var m models.Todo
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง todo ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Todo เรียงตาม created_at ล่าสุดก่อน หรือ error จาก GORM
func (r *todoRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

// ListPendingByUserID ดึง todo ที่ยัง pending ของ user
// input: userID เจ้าของข้อมูล
// output: []*Todo ที่ status=pending เรียงตาม due_at ก่อน หรือ error จาก GORM
func (r *todoRepositoryImpl) ListPendingByUserID(userID uint) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ? AND status = ?", userID, "pending").
		Order("due_at asc nulls last, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

// ListDueBetween ดึง todo ที่ due_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: []*Todo ที่ due_at อยู่ในช่วง หรือ error จาก GORM
func (r *todoRepositoryImpl) ListDueBetween(userID uint, start, end time.Time) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ? AND due_at >= ? AND due_at < ?", userID, start, end).
		Order("due_at asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

// Update บันทึกการแก้ไข todo
// input: todo domain entity ที่มี ID และค่าใหม่
// output: *Todo หลัง save หรือ error จาก GORM
func (r *todoRepositoryImpl) Update(todo *entities.Todo) (*entities.Todo, error) {
	m := models.TodoFromEntity(todo)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ todo ตาม id
// input: id ของ todo
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *todoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

// todosToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.Todo จาก GORM
// output: []*entities.Todo สำหรับส่งออกจาก repository
func todosToEntities(rows []models.Todo) []*entities.Todo {
	result := make([]*entities.Todo, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
