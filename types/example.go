package types

// Folder: types
//
// หน้าที่: เก็บ DTO / shared response type ที่ใช้ข้าม layer
//          (เช่น ResponseLogin ที่ service คืนให้ handler)
//
// แยกออกจาก:
//   - entities  : business entity ที่เป็นแก่น domain
//   - validators: request type ที่มี validate tag
//
// กติกา:
//   - เป็น struct เพียวๆ + json tag
//   - ไม่ผูกกับ framework/ORM
//   - 1 ไฟล์ = 1 กลุ่ม type (เช่น auth.go เก็บ type ตอบกลับของ auth flow)
//
// ตัวอย่าง:
//
//   package types
//
//   type ExampleResponse struct {
//       Token string `json:"token"`
//   }
