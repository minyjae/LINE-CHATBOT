package repositories

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

// conversationSessionRepositoryImpl เป็น GORM implementation ของ ConversationSessionRepository
// input: สร้างจาก NewConversationSessionRepository พร้อม *gorm.DB
// output: repository ที่อ่าน/เขียน table conversation_sessions และคืน domain entity
type conversationSessionRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ConversationSessionRepository = (*conversationSessionRepositoryImpl)(nil)

// NewConversationSessionRepository สร้าง conversation session repository
// input: db connection ของ GORM
// output: *conversationSessionRepositoryImpl ที่พร้อมใช้ query conversation session
func NewConversationSessionRepository(db *gorm.DB) *conversationSessionRepositoryImpl {
	return &conversationSessionRepositoryImpl{db: db}
}

// Create บันทึก conversation session ใหม่ลง database
// input: session domain entity ที่ต้องการบันทึก
// output: *ConversationSession ที่บันทึกแล้ว หรือ error จาก GORM
func (r *conversationSessionRepositoryImpl) Create(session *entities.ConversationSession) (*entities.ConversationSession, error) {
	m := models.ConversationSessionFromEntity(session)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetByUserIDAndKey ดึง session จาก userID และ session key
// input: userID เจ้าของ session, sessionKey ชื่อ/ประเภท session
// output: *ConversationSession ที่พบ หรือ error เช่น record not found
func (r *conversationSessionRepositoryImpl) GetByUserIDAndKey(userID uint, sessionKey string) (*entities.ConversationSession, error) {
	var m models.ConversationSession
	if err := r.db.Where("user_id = ? AND session_key = ?", userID, sessionKey).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// Upsert สร้างหรืออัปเดต session เดิมของ user/session key เดียวกัน
// input: session domain entity ที่มี userID/sessionKey/context/expiresAt
// output: *ConversationSession หลัง upsert หรือ error จาก GORM
func (r *conversationSessionRepositoryImpl) Upsert(session *entities.ConversationSession) (*entities.ConversationSession, error) {
	m := models.ConversationSessionFromEntity(session)
	if err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "session_key"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"context",
			"expires_at",
			"updated_at",
		}),
	}).Create(m).Error; err != nil {
		return nil, err
	}

	return r.GetByUserIDAndKey(m.UserID, m.SessionKey)
}

// DeleteExpired ลบ conversation session ที่หมดอายุแล้ว
// input: now เวลาปัจจุบันที่ใช้เทียบ expires_at
// output: error nil เมื่อลบสำเร็จ หรือ error จาก GORM
func (r *conversationSessionRepositoryImpl) DeleteExpired(now time.Time) error {
	return r.db.Where("expires_at <= ?", now).Delete(&models.ConversationSession{}).Error
}
