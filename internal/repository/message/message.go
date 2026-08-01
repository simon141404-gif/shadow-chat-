package message

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(ctx context.Context, msg *model.Message) error {
	query := `
		INSERT INTO messages (id, chat_id, sender_user_id, client_msg_id, message_type, content, reply_to_message_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query, msg.ID, msg.ChatID, msg.SenderUserID, msg.ClientMsgID, msg.MessageType, msg.Content, msg.ReplyToMessageID, msg.CreatedAt, msg.ExpiresAt)
	return err
}

func (r *MessageRepo) GetByID(ctx context.Context, id string) (*model.Message, error) {
	query := `
		SELECT id, chat_id, sender_user_id, client_msg_id, message_type, content, reply_to_message_id, created_at, edited_at, deleted_at, expires_at
		FROM messages WHERE id = $1
	`
	var msg model.Message
	err := r.db.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.ChatID, &msg.SenderUserID, &msg.ClientMsgID, &msg.MessageType,
		&msg.Content, &msg.ReplyToMessageID, &msg.CreatedAt, &msg.EditedAt, &msg.DeletedAt, &msg.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *MessageRepo) ListByChatID(ctx context.Context, chatID string, limit, offset int) ([]model.Message, error) {
	query := `
		SELECT id, chat_id, sender_user_id, client_msg_id, message_type, content, reply_to_message_id, created_at, edited_at, deleted_at, expires_at
		FROM messages WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderUserID, &msg.ClientMsgID, &msg.MessageType, &msg.Content, &msg.ReplyToMessageID, &msg.CreatedAt, &msg.EditedAt, &msg.DeletedAt, &msg.ExpiresAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepo) Update(ctx context.Context, msg *model.Message) error {
	query := `
		UPDATE messages SET content = $1, edited_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, msg.Content, msg.ID)
	return err
}

func (r *MessageRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE messages SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
