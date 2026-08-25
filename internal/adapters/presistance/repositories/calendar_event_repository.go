package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// calendarEventRepositoryImpl เป็น GORM implementation ของ CalendarEventRepository
// input: สร้างจาก NewCalendarEventRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table calendar_events และคืน domain entity
type calendarEventRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.CalendarEventRepository = (*calendarEventRepositoryImpl)(nil)

// NewCalendarEventRepository สร้าง calendar event repository
// input: db connection ของ GORM
// output: *calendarEventRepositoryImpl ที่พร้อมใช้ query calendar event
func NewCalendarEventRepository(db *gorm.DB) *calendarEventRepositoryImpl {
	return &calendarEventRepositoryImpl{db: db}
}

// Create บันทึก calendar event ใหม่ลง database
// input: event domain entity ที่ต้องการบันทึก
// output: *CalendarEvent ที่บันทึกแล้ว หรือ error จาก GORM
func (r *calendarEventRepositoryImpl) Create(event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	m := models.CalendarEventFromEntity(event)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง calendar event จาก primary key
// input: id ของ calendar event
// output: *CalendarEvent ที่พบ หรือ error เช่น record not found
func (r *calendarEventRepositoryImpl) GetByID(id uint) (*entities.CalendarEvent, error) {
	var m models.CalendarEvent
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง calendar event ทั้งหมดของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*CalendarEvent เรียงตาม start_at น้อยไปมาก หรือ error จาก GORM
func (r *calendarEventRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.CalendarEvent, error) {
	var rows []models.CalendarEvent
	if err := r.db.Where("user_id = ?", userID).
		Order("start_at asc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return calendarEventsToEntities(rows), nil
}

// ListByStartBetween ดึง calendar event ที่ start_at อยู่ในช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มต้นแบบ inclusive, end เวลาสิ้นสุดแบบ exclusive
// output: []*CalendarEvent ที่อยู่ในช่วงและเรียงตาม start_at หรือ error จาก GORM
func (r *calendarEventRepositoryImpl) ListByStartBetween(userID uint, start, end time.Time) ([]*entities.CalendarEvent, error) {
	var rows []models.CalendarEvent
	if err := r.db.Where("user_id = ? AND start_at >= ? AND start_at < ?", userID, start, end).
		Order("start_at asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return calendarEventsToEntities(rows), nil
}

// Update บันทึกการแก้ไข calendar event
// input: event domain entity ที่มี ID และค่าใหม่
// output: *CalendarEvent หลัง save หรือ error จาก GORM
func (r *calendarEventRepositoryImpl) Update(event *entities.CalendarEvent) (*entities.CalendarEvent, error) {
	m := models.CalendarEventFromEntity(event)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Delete ลบ calendar event ตาม id
// input: id ของ calendar event
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *calendarEventRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.CalendarEvent{}, id).Error
}

// calendarEventsToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.CalendarEvent จาก GORM
// output: []*entities.CalendarEvent สำหรับส่งออกจาก repository
func calendarEventsToEntities(rows []models.CalendarEvent) []*entities.CalendarEvent {
	result := make([]*entities.CalendarEvent, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
