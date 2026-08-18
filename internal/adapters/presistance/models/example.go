package models

// Folder: internal/adapters/presistance/models
//
// หน้าที่: เก็บ GORM model (persistence model) ที่ map ตรงกับตารางในฐานข้อมูล
//
// แยกออกจาก entities (core/domain/entities) เพื่อกัน infrastructure (GORM tag,
// column type, relation) รั่วเข้า domain layer
//
// กติกา:
//   - ใช้ struct tag ของ gorm สำหรับ schema และ json สำหรับ serialization
//   - เพิ่ม method ToEntity / FromEntity สำหรับแปลงไป-กลับกับ domain entity
//   - 1 ไฟล์ = 1 model (user.go, list.go, ...)
//
// ตัวอย่าง:
//
//   package models
//
//   import "minyjae/go-starter/internal/core/domain/entities"
//
//   type Example struct {
//       ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
//       Name      string `gorm:"type:varchar(255);not null" json:"name"`
//       CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
//   }
//
//   func (m *Example) ToEntity() *entities.Example {
//       if m == nil {
//           return nil
//       }
//       return &entities.Example{ID: m.ID, Name: m.Name}
//   }
//
//   func FromEntity(e *entities.Example) *Example {
//       if e == nil {
//           return nil
//       }
//       return &Example{ID: e.ID, Name: e.Name}
//   }
//
// อย่าลืม: register model ใน AutoMigrate ที่ internal/config/database.go
