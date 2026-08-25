package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// reminderService จัดการ business logic ของ reminder ก่อนส่งต่อไป repository
// input: สร้างจาก NewReminderServiceImpl พร้อม ReminderRepository
// output: service ที่ทำ create/list/update/delete reminder โดยผูกข้อมูลกับ userID
type reminderService struct {
	repo repoPort.ReminderRepository
}

var _ servicePort.ReminderService = (*reminderService)(nil)

// NewReminderServiceImpl สร้าง reminder service implementation
// input: repo repository สำหรับอ่าน/เขียน reminder
// output: *reminderService ที่พร้อมถูกใช้ผ่าน ReminderService interface
func NewReminderServiceImpl(repo repoPort.ReminderRepository) *reminderService {
	return &reminderService{repo: repo}
}

// Create สร้าง reminder ใหม่ให้ user และเติมค่า default ที่จำเป็น
// input: userID เจ้าของ reminder, reminder ข้อมูลเตือนความจำที่รับจาก controller/assistant
// output: *Reminder ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
func (s *reminderService) Create(userID uint, reminder *entities.Reminder) (*entities.Reminder, error) {
	now := time.Now()
	reminder.UserID = userID
	reminder.Status = defaultString(reminder.Status, "pending")
	reminder.CreatedAt = now
	reminder.UpdatedAt = now
	return s.repo.Create(reminder)
}

// List ดึง reminder ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Reminder รายการ reminder หรือ error จาก repository
func (s *reminderService) List(userID uint, limit, offset int) ([]*entities.Reminder, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// Update แก้ไข reminder โดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id reminder ที่ต้องการแก้, reminder ค่าใหม่
// output: *Reminder ที่แก้แล้ว หรือ ErrForbidden ถ้า reminder ไม่ใช่ของ user นี้
func (s *reminderService) Update(userID, id uint, reminder *entities.Reminder) (*entities.Reminder, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	reminder.ID = id
	reminder.UserID = userID
	reminder.CreatedAt = current.CreatedAt
	reminder.UpdatedAt = time.Now()
	if reminder.Status == "" {
		reminder.Status = current.Status
	}
	return s.repo.Update(reminder)
}

// Delete ลบ reminder โดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id reminder ที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้า reminder ไม่ใช่ของ user นี้
func (s *reminderService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
