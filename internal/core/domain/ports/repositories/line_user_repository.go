package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// LineUserRepository คือ contract สำหรับอ่าน/เขียนบัญชี LINE ที่ผูกกับ user
// input: line user entity, id, lineUserID หรือ userID
// output: line user entity/list หรือ error จาก adapter ที่ implement จริง
type LineUserRepository interface {
	Create(lineUser *entities.LineUser) (*entities.LineUser, error)
	GetByID(id uint) (*entities.LineUser, error)
	GetByLineUserID(lineUserID string) (*entities.LineUser, error)
	ListByUserID(userID uint) ([]*entities.LineUser, error)
	Update(lineUser *entities.LineUser) (*entities.LineUser, error)
	Delete(id uint) error
}
