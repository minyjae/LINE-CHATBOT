package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// reminderRepositoryImpl เป็น GORM implementation ของ ReminderRepository
// input: สร้างจาก NewReminderRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table reminders และคืน domain entity
type reminderRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ReminderRepository = (*reminderRepositoryImpl)(nil)

// NewReminderRepository สร้าง reminder repository
// input: db connection ของ GORM
// output: *reminderRepositoryImpl ที่พร้อมใช้ query reminder
func NewReminderRepository(db *gorm.DB) *reminderRepositoryImpl {
	return &reminderRepositoryImpl{db: db}
}

// Create บันทึก reminder ใหม่ลง database
// input: reminder domain entity ที่ต้องการบันทึก
// output: *Reminder ที่บันทึกแล้ว หรือ error จาก GORM
func (r *reminderRepositoryImpl) Create(reminder *entities.Reminder) (*entities.Reminder, error) {
	m := models.ReminderFromEntity(reminder)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByID ดึง reminder จาก primary key
// input: id ของ reminder
// output: *Reminder ที่พบ หรือ error เช่น record not found
func (r *reminderRepositoryImpl) GetByID(id uint) (*entities.Reminder, error) {
	var m models.Reminder
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// ListByUserID ดึง reminder ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Reminder เรียงตาม remind_at น้อยไปมาก หรือ error จาก GORM
func (r *reminderRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Reminder, error) {
	var rows []models.Reminder
	if err := r.db.Where("user_id = ?", userID).
		Order("remind_at asc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return remindersToEntities(rows), nil
}

// ListPendingDue ดึง reminder ที่ถึงเวลาต้องเตือนแล้ว
// input: now เวลาปัจจุบัน, limit จำนวนสูงสุด
// output: []*Reminder ที่ status=pending และ remind_at <= now หรือ error จาก GORM
func (r *reminderRepositoryImpl) ListPendingDue(now time.Time, limit int) ([]*entities.Reminder, error) {
	var rows []models.Reminder
	if err := r.db.Where("status = ? AND remind_at <= ?", "pending", now).
		Order("remind_at asc").
		Limit(normalizeLimit(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return remindersToEntities(rows), nil
}

// ListPendingDueOrUpcoming ดึง reminder สำหรับ worker ทั้งเตือนล่วงหน้าและเตือนจริง
// input: now เวลาปัจจุบัน, upcomingUntil เวลาปลายทางสำหรับ pre-alert, limit จำนวนสูงสุด
// output: []*Reminder ที่ยัง pending และใกล้ถึงเวลา หรือ pre_sent ที่ถึงเวลาจริงแล้ว
func (r *reminderRepositoryImpl) ListPendingDueOrUpcoming(now, upcomingUntil time.Time, limit int) ([]*entities.Reminder, error) {
	var rows []models.Reminder
	if err := r.db.Where(
		"(status = ? AND remind_at <= ?) OR (status = ? AND remind_at <= ?)",
		"pending", upcomingUntil,
		"pre_sent", now,
	).
		Order("remind_at asc").
		Limit(normalizeLimit(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return remindersToEntities(rows), nil
}

// Update บันทึกการแก้ไข reminder
// input: reminder domain entity ที่มี ID และค่าใหม่
// output: *Reminder หลัง save หรือ error จาก GORM
func (r *reminderRepositoryImpl) Update(reminder *entities.Reminder) (*entities.Reminder, error) {
	m := models.ReminderFromEntity(reminder)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// MarkPreSent เปลี่ยน reminder เป็นสถานะ pre_sent หลังส่งเตือนล่วงหน้าแล้ว
// input: id ของ reminder, sentAt เวลาที่ส่ง pre-alert
// output: error nil เมื่อ update สำเร็จ หรือ error จาก GORM
func (r *reminderRepositoryImpl) MarkPreSent(id uint, sentAt time.Time) error {
	return r.db.Model(&models.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":  "pre_sent",
			"sent_at": sentAt,
		}).Error
}

// MarkSent เปลี่ยน reminder เป็นสถานะ sent หลังส่งเตือนถึงเวลาจริงแล้ว
// input: id ของ reminder, sentAt เวลาที่ส่ง reminder จริง
// output: error nil เมื่อ update สำเร็จ หรือ error จาก GORM
func (r *reminderRepositoryImpl) MarkSent(id uint, sentAt time.Time) error {
	return r.db.Model(&models.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":  "sent",
			"sent_at": sentAt,
		}).Error
}

// Delete ลบ reminder ตาม id
// input: id ของ reminder
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *reminderRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Reminder{}, id).Error
}

// remindersToEntities แปลง slice ของ persistence model เป็น domain entity
// input: rows []models.Reminder จาก GORM
// output: []*entities.Reminder สำหรับส่งออกจาก repository
func remindersToEntities(rows []models.Reminder) []*entities.Reminder {
	result := make([]*entities.Reminder, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
