package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// IntegrationAccountRepository คือ contract สำหรับอ่าน/เขียนบัญชี integration ของ user
// input: integration account entity, id, userID หรือ provider
// output: integration account entity/list หรือ error จาก adapter ที่ implement จริง
type IntegrationAccountRepository interface {
	Create(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error)
	GetByID(id uint) (*entities.IntegrationAccount, error)
	GetByUserIDAndProvider(userID uint, provider string) (*entities.IntegrationAccount, error)
	ListByUserID(userID uint) ([]*entities.IntegrationAccount, error)
	Update(account *entities.IntegrationAccount) (*entities.IntegrationAccount, error)
	Delete(id uint) error
}
