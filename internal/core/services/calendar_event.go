package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// calendarEventService จัดการ business logic ของ calendar event ก่อนส่งต่อไป repository
// input: สร้างจาก NewCalendarEventServiceImpl พร้อม CalendarEventRepository
// output: service ที่ทำ create/list/list-by-date/update/delete calendar event โดยผูกข้อมูลกับ userID
type calendarEventService struct {
	repo repoPort.CalendarEventRepository
}

var _ servicePort.CalendarEventService = (*calendarEventService)(nil)

// NewCalendarEventServiceImpl สร้าง calendar event service implementation
// input: repo repository สำหรับอ่าน/เขียน calendar event
// output: *calendarEventService ที่พร้อมถูกใช้ผ่าน CalendarEventService interface
func NewCalendarEventServiceImpl(repo repoPort.CalendarEventRepository) *calendarEventService {
	return &calendarEventService{repo: repo}
}

// Create สร้าง calendar event ใหม่ให้ user และเติมค่า default ที่จำเป็น
// input: userID เจ้าของนัด, event ข้อมูลนัดจาก controller/assistant
// output: *CalendarEvent ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
func (s *calendarEventService) Create(userID uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	now := time.Now()
	event.UserID = userID
	event.SyncStatus = defaultString(event.SyncStatus, "local")
	event.CreatedAt = now
	event.UpdatedAt = now
	return s.repo.Create(event)
}

// List ดึง calendar event ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*CalendarEvent รายการนัด หรือ error จาก repository
func (s *calendarEventService) List(userID uint, limit, offset int) ([]*entities.CalendarEvent, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// ListByDate ดึง calendar event ของ user เฉพาะวันเดียวตาม timezone
// input: userID เจ้าของข้อมูล, date วันที่ต้องการดู, loc timezone สำหรับคำนวณต้น/ท้ายวัน
// output: []*CalendarEvent ที่ StartAt อยู่ในวันนั้น หรือ error จาก repository
func (s *calendarEventService) ListByDate(userID uint, date time.Time, loc *time.Location) ([]*entities.CalendarEvent, error) {
	if loc == nil {
		loc = time.Local
	}
	localDate := date.In(loc)
	start := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)
	return s.repo.ListByStartBetween(userID, start, end)
}

// Update แก้ไข calendar event โดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id event ที่ต้องการแก้, event ค่าใหม่
// output: *CalendarEvent ที่แก้แล้ว หรือ ErrForbidden ถ้า event ไม่ใช่ของ user นี้
func (s *calendarEventService) Update(userID, id uint, event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	event.ID = id
	event.UserID = userID
	event.CreatedAt = current.CreatedAt
	event.UpdatedAt = time.Now()
	if event.SyncStatus == "" {
		event.SyncStatus = current.SyncStatus
	}
	return s.repo.Update(event)
}

// Delete ลบ calendar event โดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id event ที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้า event ไม่ใช่ของ user นี้
func (s *calendarEventService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
