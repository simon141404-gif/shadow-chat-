package repository

import (
	"github.com/yourorg/shadowchat/backend/internal/repository/attachment"
	"github.com/yourorg/shadowchat/backend/internal/repository/chat"
	"github.com/yourorg/shadowchat/backend/internal/repository/contact"
	"github.com/yourorg/shadowchat/backend/internal/repository/message"
	"github.com/yourorg/shadowchat/backend/internal/repository/notification"
	"github.com/yourorg/shadowchat/backend/internal/repository/session"
	"github.com/yourorg/shadowchat/backend/internal/repository/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Bundle struct {
	Users        *user.UserRepo
	Chats        *chat.ChatRepo
	Messages     *message.MessageRepo
	Attachments  *attachment.AttachmentRepo
	Sessions     *session.SessionRepo
	Contacts     *contact.ContactRepo
	Notifications *notification.NotificationRepo
}

func NewBundle(pg *pgxpool.Pool, redis *redis.Client) Bundle {
	return Bundle{
		Users:        user.NewUserRepo(pg),
		Chats:        chat.NewChatRepo(pg),
		Messages:     message.NewMessageRepo(pg),
		Attachments:  attachment.NewAttachmentRepo(pg),
		Sessions:     session.NewSessionRepo(pg, redis),
		Contacts:     contact.NewContactRepo(pg),
		Notifications: notification.NewNotificationRepo(pg, redis),
	}
}
