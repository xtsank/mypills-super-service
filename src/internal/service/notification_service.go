package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/dto"
)

type INotificationService interface {
	SendDueNotifications(ctx context.Context)
}

type NotificationService struct {
	userRepo    user.IUserRepository
	cabinetRepo cabinet_item.ICabinetItemRepository
	emailSender EmailSender
	logger      *slog.Logger
}

func NewNotificationService(i do.Injector) (INotificationService, error) {
	userRepo := do.MustInvoke[user.IUserRepository](i)
	cabinetRepo := do.MustInvoke[cabinet_item.ICabinetItemRepository](i)
	emailSender := do.MustInvoke[EmailSender](i)
	logger := do.MustInvoke[*slog.Logger](i)

	return &NotificationService{
		userRepo:    userRepo,
		cabinetRepo: cabinetRepo,
		emailSender: emailSender,
		logger:      logger,
	}, nil
}

func (s *NotificationService) SendDueNotifications(ctx context.Context) {
	users, err := s.userRepo.FindNotifyEnabled(ctx)
	if err != nil {
		s.logger.Error("notify_list_users_failed", slog.Any("error", err))
		return
	}

	now := time.Now()
	for _, u := range users {
		if u.Email == "" || u.Notify == nil || !u.Notify.Enabled {
			continue
		}
		if !isDue(u.Notify.LastNotifiedAt, u.Notify.IntervalMinutes, now) {
			continue
		}

		items, err := s.cabinetRepo.FindExpiredByUserID(ctx, u.ID)
		if err != nil {
			s.logger.Error("notify_expired_fetch_failed", slog.Any("error", err), slog.String("user_id", u.ID.String()))
			continue
		}
		if len(items) == 0 {
			continue
		}

		subject := "Просроченные лекарства"
		body := buildExpiredEmailBody(items)
		
		if err := s.emailSender.Send(ctx, u.Email, subject, body); err != nil {
			s.logger.Error("notify_email_failed", slog.Any("error", err), slog.String("user_id", u.ID.String()))
			continue
		}

		if err := s.userRepo.UpdateNotify(ctx, u.ID, u.Notify.Enabled, u.Notify.IntervalMinutes, &now); err != nil {
			s.logger.Error("notify_update_failed", slog.Any("error", err), slog.String("user_id", u.ID.String()))
		}
	}
}

func isDue(last *time.Time, intervalMinutes int, now time.Time) bool {
	if intervalMinutes <= 0 {
		return false
	}
	if last == nil {
		return true
	}
	return last.Add(time.Duration(intervalMinutes)*time.Minute).Before(now) || last.Add(time.Duration(intervalMinutes)*time.Minute).Equal(now)
}

func buildExpiredEmailBody(items []*dto.ExpiredItemDto) string {
	var b strings.Builder
	b.WriteString("У вас есть просроченные лекарства:\n\n")
	for _, item := range items {
		line := fmt.Sprintf("- %s, изготовлено: %s, годно до: %s, количество: %.2f\n",
			item.MedicineName,
			item.DateOfManufacture.Format("2006-01-02"),
			item.ExpiresAt.Format("2006-01-02"),
			item.Quantity,
		)
		b.WriteString(line)
	}
	return b.String()
}
