package services

import "errors"

// Folder: internal/core/services
//
// หน้าที่: implement port ของ service (internal/core/domain/ports/services)
//          เป็นที่อยู่ของ business logic เพียวๆ (validation, orchestration,
//          แปลง entity, จัดการ error)
//
// กฎสำคัญ:
//   - รับ dependency เป็น "port repository" (interface) เท่านั้น
//     ห้ามยึดติดกับ concrete persistence
//   - แปลงข้อผิดพลาดจากชั้นล่าง (เช่น gorm.ErrRecordNotFound) เป็น
//     domain error ที่ประกาศไว้ใน errors.go ของ package นี้
//   - 1 ไฟล์ = 1 service implementation
//   - errors.go : ประกาศ sentinel error ทั้งหมดของ service layer ไว้รวมกัน
//
// ตัวอย่าง sentinel error (ย้ายไป errors.go จริงเมื่อเริ่มใช้งาน):
var ErrExampleNotFound = errors.New("example not found")

// ตัวอย่าง implementation:
//
//   package services
//
//   import (
//       "errors"
//
//       "minyjae/go-starter/internal/core/domain/entities"
//       repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
//       "gorm.io/gorm"
//   )
//
//   type exampleService struct {
//       repo repoPort.ExampleRepository
//   }
//
//   func NewExampleServiceImpl(r repoPort.ExampleRepository) *exampleService {
//       return &exampleService{repo: r}
//   }
//
//   func (s *exampleService) Get(id uint) (*entities.Example, error) {
//       e, err := s.repo.GetByID(id)
//       if err != nil {
//           if errors.Is(err, gorm.ErrRecordNotFound) {
//               return nil, ErrExampleNotFound
//           }
//           return nil, err
//       }
//       return e, nil
//   }
