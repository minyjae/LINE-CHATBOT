package entities

// Folder: internal/core/domain/entities
//
// หน้าที่: เก็บ "domain entity" ที่เป็นแก่นของ business เพียวๆ
//
// กฎสำคัญ:
//   - ห้าม import infrastructure (gorm, fiber, sqlx, ฯลฯ) เด็ดขาด
//   - เป็น plain Go struct + business method เท่านั้น
//   - ใช้ json tag ได้ (เพราะ entity ถูกส่งออกผ่าน handler เป็น response)
//   - 1 ไฟล์ = 1 entity (user.go, list.go, ...)
//
// ตัวอย่าง:
//
//   package entities
//
//   type Example struct {
//       ID   uint   `json:"id"`
//       Name string `json:"name"`
//   }
//
//   // business method สามารถเขียนไว้ใน entity ได้
//   func (e *Example) IsValid() bool {
//       return e.Name != ""
//   }
