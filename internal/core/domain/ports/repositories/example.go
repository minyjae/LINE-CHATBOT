package repositories

// Folder: internal/core/domain/ports/repositories
//
// หน้าที่: ประกาศ "port" ของ repository (interface) ที่ domain ต้องการ
//          ส่วน implement อยู่ที่ internal/adapters/presistance/repositories
//
// กฎสำคัญ (Hexagonal / Clean Architecture):
//   - import ได้แค่ entities, validators หรือ standard library
//   - ห้าม import gorm / fiber / persistence model
//   - 1 ไฟล์ = 1 repository interface (auth_repository.go, ...)
//
// ตัวอย่าง:
//
//   package repositories
//
//   import "minyjae/go-starter/internal/core/domain/entities"
//
//   type ExampleRepository interface {
//       GetByID(id uint) (*entities.Example, error)
//       Create(e entities.Example) (*entities.Example, error)
//   }
