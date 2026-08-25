package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/types"
)

// lineWebhookService จัดการ use-case ข้อความจาก LINE ก่อนส่งเข้า assistant
// input: สร้างจาก NewLineWebhookServiceImpl พร้อม user/line user/message log repository และ assistant service
// output: service ที่รับ LINE text, ensure user, บันทึก log, และคืน reply text
type lineWebhookService struct {
	userRepo       repoPort.UserRepository
	lineUserRepo   repoPort.LineUserRepository
	messageLogRepo repoPort.MessageLogRepository
	assistant      servicePort.AssistantService
}

var _ servicePort.LineWebhookService = (*lineWebhookService)(nil)

// NewLineWebhookServiceImpl สร้าง LINE webhook service implementation
// input: repository สำหรับ user/line user/message log และ assistant service
// output: *lineWebhookService ที่พร้อมถูกใช้ผ่าน LineWebhookService interface
func NewLineWebhookServiceImpl(
	userRepo repoPort.UserRepository,
	lineUserRepo repoPort.LineUserRepository,
	messageLogRepo repoPort.MessageLogRepository,
	assistant servicePort.AssistantService,
) *lineWebhookService {
	return &lineWebhookService{
		userRepo:       userRepo,
		lineUserRepo:   lineUserRepo,
		messageLogRepo: messageLogRepo,
		assistant:      assistant,
	}
}

// HandleTextMessage รับข้อความ LINE หนึ่งข้อความแล้วสร้างผลลัพธ์ตอบกลับ
// input: LineTextMessageInput ที่มี line user id, reply token, text และ now
// output: LineTextMessageResult ที่มี user id, message log id, intent, reply text หรือ error
func (s *lineWebhookService) HandleTextMessage(input types.LineTextMessageInput) (*types.LineTextMessageResult, error) {
	now := input.Now
	if now.IsZero() {
		now = time.Now()
	}

	user, err := s.ensureUser(input.LineUserID, now)
	if err != nil {
		return nil, err
	}

	inbound, err := s.messageLogRepo.Create(&entities.MessageLog{
		UserID:         user.ID,
		Source:         "line",
		Direction:      "inbound",
		MessageType:    "text",
		MessageText:    input.Text,
		LineReplyToken: &input.ReplyToken,
		CreatedAt:      now,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.assistant.HandleTextMessage(types.AssistantMessageInput{
		UserID:       user.ID,
		MessageLogID: &inbound.ID,
		Text:         input.Text,
		Now:          now,
		Timezone:     user.Timezone,
	})
	if err != nil {
		return nil, err
	}

	if _, err := s.messageLogRepo.Create(&entities.MessageLog{
		UserID:      user.ID,
		Source:      "line",
		Direction:   "outbound",
		MessageType: "text",
		MessageText: result.ReplyText,
		CreatedAt:   now,
	}); err != nil {
		return nil, err
	}

	return &types.LineTextMessageResult{
		UserID:       user.ID,
		MessageLogID: inbound.ID,
		Intent:       result.Intent,
		ReplyText:    result.ReplyText,
	}, nil
}

// ensureUser หา user จาก LINE user id หรือสร้าง user ใหม่ถ้ายังไม่เคยเจอ
// input: lineUserID จาก LINE, now เวลาปัจจุบันสำหรับ LastSeenAt/timestamps
// output: *User ของระบบที่ผูกกับ LINE user id หรือ error จาก repository
func (s *lineWebhookService) ensureUser(lineUserID string, now time.Time) (*entities.User, error) {
	lineUser, err := s.lineUserRepo.GetByLineUserID(lineUserID)
	if err == nil {
		user, err := s.userRepo.GetByID(lineUser.UserID)
		if err != nil {
			return nil, err
		}

		user.LastSeenAt = &now
		if _, err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		lineUser.LastSeenAt = &now
		if _, err := s.lineUserRepo.Update(lineUser); err != nil {
			return nil, err
		}
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user, err := s.userRepo.Create(&entities.User{
		DisplayName: "LINE User",
		Timezone:    "Asia/Bangkok",
		Locale:      "th-TH",
		Status:      "active",
		LastSeenAt:  &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}

	if _, err := s.lineUserRepo.Create(&entities.LineUser{
		UserID:     user.ID,
		LineUserID: lineUserID,
		Status:     "active",
		LastSeenAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, err
	}

	return user, nil
}
