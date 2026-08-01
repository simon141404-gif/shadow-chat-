package notification

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type NotificationRepo struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewNotificationRepo(db *pgxpool.Pool, redis *redis.Client) *NotificationRepo {
	return &NotificationRepo{db: db, redis: redis}
}

func (r *NotificationRepo) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, body, data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, notification.ID, notification.UserID, notification.Type, notification.Title, notification.Body, notification.Data)
	return err
}

func (r *NotificationRepo) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, body, data, read_at, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Data, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications SET read_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *NotificationRepo) RegisterPushToken(ctx context.Context, token *model.PushToken) error {
	query := `
		INSERT INTO push_tokens (id, user_id, token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, token) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, token.ID, token.UserID, token.Token)
	return err
}

func (r *NotificationRepo) GetPushTokens(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT token FROM push_tokens WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// PubSub publishes notification to user
func (r *NotificationRepo) PublishNotification(ctx context.Context, userID, message string) error {
	return r.redis.Publish(ctx, "notifications:"+userID, message).Err()
}

// Subscribe returns a channel for user notifications
func (r *NotificationRepo) Subscribe(ctx context.Context, userID string) *redis.PubSub {
	return r.redis.Subscribe(ctx, "notifications:"+userID)
}
