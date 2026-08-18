package repositories

// Folder: internal/adapters/presistance/repositories
//
// หน้าที่: implement port (interface) ของ repository ที่ประกาศไว้ใน
//          internal/core/domain/ports/repositories โดยใช้ GORM ติดต่อ DB
//
// กติกา:
//   - 1 ไฟล์ = 1 repository (auth_repository.go, list_repository.go, ...)
//   - struct ตั้งชื่อ <name>RepositoryImpl, ขึ้นต้นด้วยตัวเล็ก (unexported)
//   - constructor NewXxxRepository(db *gorm.DB) คืน pointer ของ struct
//   - แปลง persistence model → domain entity ก่อนส่งออก (ห้ามให้ models รั่ว
//     ออกนอก folder นี้ ยกเว้นเคสที่ port กำหนดให้คืน *models.X เป็นพิเศษ)
//
// ตัวอย่าง:
//
//   package repositories
//
//   import (
//       "minyjae/go-starter/internal/adapters/presistance/models"
//       "minyjae/go-starter/internal/core/domain/entities"
//       "gorm.io/gorm"
//   )
//
//   type exampleRepositoryImpl struct {
//       db *gorm.DB
//   }
//
//   func NewExampleRepository(db *gorm.DB) *exampleRepositoryImpl {
//       return &exampleRepositoryImpl{db: db}
//   }
//
//   func (r *exampleRepositoryImpl) GetByID(id uint) (*entities.Example, error) {
//       var m models.Example
//       if err := r.db.First(&m, id).Error; err != nil {
//           return nil, err
//       }
//       return m.ToEntity(), nil
//   }
