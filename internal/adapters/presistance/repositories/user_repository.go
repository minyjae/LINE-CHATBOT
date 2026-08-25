package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// userRepositoryImpl เป็น GORM implementation ของ UserRepository
// input: สร้างจาก NewUserRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table users และคืน domain entity
type userRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.UserRepository = (*userRepositoryImpl)(nil)

// NewUserRepository สร้าง user repository
// input: db connection ของ GORM
// output: *userRepositoryImpl ที่พร้อมใช้ query user
func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

// Create บันทึก user ใหม่ลง database
// input: user domain entity ที่ต้องการบันทึก
// output: *User ที่บันทึกแล้ว หรือ error จาก GORM
func (r *userRepositoryImpl) Create(user *entities.User) (*entities.User, error) {
	m := models.FromEntity(user)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง user จาก primary key
// input: id ของ user
// output: *User ที่พบ หรือ error เช่น record not found
func (r *userRepositoryImpl) GetByID(id uint) (*entities.User, error) {
	var m models.User
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByEmail ดึง user จาก email
// input: email ที่ต้องการค้นหา
// output: *User ที่พบ หรือ error เช่น record not found
func (r *userRepositoryImpl) GetByEmail(email string) (*entities.User, error) {
	var m models.User
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Update บันทึกการแก้ไข user
// input: user domain entity ที่มี ID และค่าใหม่
// output: *User หลัง save หรือ error จาก GORM
func (r *userRepositoryImpl) Update(user *entities.User) (*entities.User, error) {
	m := models.FromEntity(user)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ user ตาม id
// input: id ของ user
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
