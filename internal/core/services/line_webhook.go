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

type lineWebhookService struct {
	userRepo       repoPort.UserRepository
	lineUserRepo   repoPort.LineUserRepository
	messageLogRepo repoPort.MessageLogRepository
	assistant      servicePort.AssistantService
}

var _ servicePort.LineWebhookService = (*lineWebhookService)(nil)

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
