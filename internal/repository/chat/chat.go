package chat

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo struct {
	db *pgxpool.Pool
}

func NewChatRepo(db *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) Create(ctx context.Context, chat *model.Chat) error {
	query := `
		INSERT INTO chats (id, type, name, avatar_url)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, chat.ID, chat.Type, chat.Name, chat.AvatarURL)
	return err
}

func (r *ChatRepo) GetByID(ctx context.Context, id string) (*model.Chat, error) {
	query := `
		SELECT id, type, name, avatar_url, created_at, updated_at, deleted_at
		FROM chats WHERE id = $1 AND deleted_at IS NULL
	`
	var chat model.Chat
	err := r.db.QueryRow(ctx, query, id).Scan(
		&chat.ID, &chat.Type, &chat.Name, &chat.AvatarURL,
		&chat.CreatedAt, &chat.UpdatedAt, &chat.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepo) ListByUserID(ctx context.Context, userID string) ([]model.Chat, error) {
	query := `
		SELECT c.id, c.type, c.name, c.avatar_url, c.created_at, c.updated_at, c.deleted_at
		FROM chats c
		JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.updated_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.Name, &chat.AvatarURL, &chat.CreatedAt, &chat.UpdatedAt, &chat.DeletedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepo) AddMember(ctx context.Context, member *model.ChatMember) error {
	query := `
		INSERT INTO chat_members (chat_id, user_id, role, joined_at, invited_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, member.ChatID, member.UserID, member.Role, member.JoinedAt, member.InvitedBy)
	return err
}

func (r *ChatRepo) GetMembers(ctx context.Context, chatID string) ([]model.ChatMember, error) {
	query := `
		SELECT chat_id, user_id, role, joined_at, invited_by
		FROM chat_members WHERE chat_id = $1
	`
	rows, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.ChatMember
	for rows.Next() {
		var m model.ChatMember
		if err := rows.Scan(&m.ChatID, &m.UserID, &m.Role, &m.JoinedAt, &m.InvitedBy); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *ChatRepo) IsMember(ctx context.Context, chatID, userID string) (bool, error) {
	query := `SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2`
	var exists int
	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return false, nil
	}
	return true, nil
}
