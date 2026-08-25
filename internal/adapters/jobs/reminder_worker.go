package jobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	lineAdapter "minyjae/go-starter/internal/adapters/line"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type ReminderWorker struct {
	reminderRepo   repoPort.ReminderRepository
	lineUserRepo   repoPort.LineUserRepository
	messageLogRepo repoPort.MessageLogRepository
	messenger      lineAdapter.Messenger
	interval       time.Duration
	batchSize      int
}

func NewReminderWorker(
	reminderRepo repoPort.ReminderRepository,
	lineUserRepo repoPort.LineUserRepository,
	messageLogRepo repoPort.MessageLogRepository,
	messenger lineAdapter.Messenger,
	interval time.Duration,
	batchSize int,
) *ReminderWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	if batchSize <= 0 {
		batchSize = 50
	}

	return &ReminderWorker{
		reminderRepo:   reminderRepo,
		lineUserRepo:   lineUserRepo,
		messageLogRepo: messageLogRepo,
		messenger:      messenger,
		interval:       interval,
		batchSize:      batchSize,
	}
}

func (w *ReminderWorker) Start(ctx context.Context) {
	go func() {
		log.Printf("Reminder worker started interval=%s batch_size=%d", w.interval, w.batchSize)
		w.processDueReminders()

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Reminder worker stopped")
				return
			case <-ticker.C:
				w.processDueReminders()
			}
		}
	}()
}

func (w *ReminderWorker) processDueReminders() {
	now := time.Now()
	reminders, err := w.reminderRepo.ListPendingDue(now, w.batchSize)
	if err != nil {
		log.Printf("Reminder worker failed to list reminders: %v", err)
		return
	}

	for _, reminder := range reminders {
		if err := w.sendReminder(reminder, now); err != nil {
			log.Printf("Reminder worker failed to send reminder id=%d user_id=%d: %v", reminder.ID, reminder.UserID, err)
			continue
		}
	}
}

func (w *ReminderWorker) sendReminder(reminder *entities.Reminder, sentAt time.Time) error {
	lineUsers, err := w.lineUserRepo.ListByUserID(reminder.UserID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("ถึงเวลาแล้วค่ะ เลขามาสะกิดเรื่อง \"%s\" ค่ะ", reminder.Title)
	sentCount := 0
	for _, lineUser := range lineUsers {
		if lineUser.Status != "active" {
			continue
		}

		if err := w.messenger.PushText(lineUser.LineUserID, message); err != nil {
			if errors.Is(err, lineAdapter.ErrMessengerNotConfigured) {
				return err
			}
			log.Printf("Reminder worker push failed reminder_id=%d line_user_id=%s: %v", reminder.ID, lineUser.LineUserID, err)
			continue
		}
		sentCount++
	}

	if sentCount == 0 {
		return fmt.Errorf("no active LINE recipient for reminder")
	}

	if _, err := w.messageLogRepo.Create(&entities.MessageLog{
		UserID:      reminder.UserID,
		Source:      "line",
		Direction:   "outbound",
		MessageType: "text",
		MessageText: message,
		CreatedAt:   sentAt,
	}); err != nil {
		return err
	}

	return w.reminderRepo.MarkSent(reminder.ID, sentAt)
}
