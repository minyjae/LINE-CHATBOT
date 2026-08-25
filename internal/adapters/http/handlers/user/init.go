package user

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

// IUserController กำหนด contract ของ user controller
// input: ยังไม่มี method ที่เปิดใช้ในตอนนี้
// output: interface สำหรับเตรียมผูก user handler ในอนาคต
type IUserController interface {
	// GetUser(c *fiber.Ctx) error
}

// userController เก็บ GenericController สำหรับ user handler
// input: สร้างจาก NewUserController
// output: controller ที่พร้อมเพิ่ม user action ภายหลัง
type userController struct {
	util_controllers.GenericController
	// userService services.UserService
}

// NewUserController สร้าง user controller
// input: ไม่มี dependency ในตอนนี้
// output: IUserController ที่พร้อมผูก route เมื่อมี user action
func NewUserController() IUserController {
	return &userController{
		GenericController: util_controllers.NewGenericController(),
		// userService:       userService,
	}
}
