package services

// Folder: internal/core/domain/ports/services
//
// หน้าที่: ประกาศ "port" ของ service (interface) ที่ adapter ฝั่ง HTTP
//          (handler) จะเรียกใช้ ส่วน implement อยู่ที่ internal/core/services
//
// กฎสำคัญ:
//   - import ได้แค่ entities, types, validators หรือ standard library
//   - ห้าม import infrastructure (gorm/fiber/...)
//   - 1 ไฟล์ = 1 service interface
//
// ตัวอย่าง:
//
//   package services
//
//   import "minyjae/go-starter/internal/core/domain/entities"
//
//   type ExampleService interface {
//       Get(id uint) (*entities.Example, error)
//       Create(e entities.Example) (*entities.Example, error)
//   }
