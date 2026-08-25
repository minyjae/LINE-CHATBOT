package auth

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

// IAuthController กำหนด contract ของ auth controller
// input: ไม่มี method ในตอนนี้
// output: interface ว่างสำหรับเตรียมผูก auth handler ในอนาคต
type IAuthController interface {
}

// authController เก็บ GenericController สำหรับ auth handler
// input: สร้างจาก NewAuthController
// output: controller ที่พร้อมเพิ่ม auth action ภายหลัง
type authController struct {
	util_controllers.GenericController
}

// NewAuthController สร้าง auth controller
// input: ไม่มี dependency ในตอนนี้
// output: IAuthController ที่พร้อมผูก route เมื่อมี auth action
func NewAuthController() IAuthController {
	return &authController{
		GenericController: util_controllers.NewGenericController(),
	}
}
