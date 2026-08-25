package list

import (
	util_controllers "minyjae/go-starter/utils/controller"
)

// IListController กำหนด contract ของ list controller
// input: ไม่มี method ในตอนนี้
// output: interface ว่างสำหรับเตรียมผูก list handler ในอนาคต
type IListController interface {
}

// listController เก็บ GenericController สำหรับ list handler
// input: สร้างจาก NewListController
// output: controller ที่พร้อมเพิ่ม list action ภายหลัง
type listController struct {
	util_controllers.GenericController
}

// NewListController สร้าง list controller
// input: ไม่มี dependency ในตอนนี้
// output: IListController ที่พร้อมผูก route เมื่อมี list action
func NewListController() IListController {
	return &listController{
		GenericController: util_controllers.NewGenericController(),
	}
}
