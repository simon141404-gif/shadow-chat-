package notifications

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/notification"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationService struct {
	notifications *notification.NotificationRepo
	logger        *zap.Logger
}

func NewNotificationService(notifications *notification.NotificationRepo, logger *zap.Logger) *NotificationService {
	return &NotificationService{notifications: notifications, logger: logger}
}

type RegisterPushTokenRequest struct {
	Token string `json:"token"`
}

func (s *NotificationService) RegisterPushToken(ctx context.Context, userID string, req RegisterPushTokenRequest) error {
	token := &model.PushToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     req.Token,
		CreatedAt: time.Now(),
	}
	return s.notifications.RegisterPushToken(ctx, token)
}

func (s *NotificationService) List(ctx context.Context, userID string, limit, offset int) ([]model.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notifications.ListByUserID(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	return s.notifications.MarkAsRead(ctx, notificationID)
}

func (s *NotificationService) SendNotification(ctx context.Context, userID, title, body string, data string) error {
	notification := &model.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      model.NotificationTypeMessage,
		Title:     title,
		Body:      body,
		Data:      data,
		CreatedAt: time.Now(),
	}

	if err := s.notifications.Create(ctx, notification); err != nil {
		return err
	}

	// Publish to Redis for real-time delivery
	return s.notifications.PublishNotification(ctx, userID, notification.ID)
}
