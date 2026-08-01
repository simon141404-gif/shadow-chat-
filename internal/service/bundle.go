package service

import (
	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/repository"
	"github.com/yourorg/shadowchat/backend/internal/service/auth"
	"github.com/yourorg/shadowchat/backend/internal/service/chats"
	"github.com/yourorg/shadowchat/backend/internal/service/contacts"
	"github.com/yourorg/shadowchat/backend/internal/service/groups"
	"github.com/yourorg/shadowchat/backend/internal/service/messages"
	"github.com/yourorg/shadowchat/backend/internal/service/notifications"
	"github.com/yourorg/shadowchat/backend/internal/service/profile"
	"github.com/yourorg/shadowchat/backend/internal/service/uploads"
	"github.com/yourorg/shadowchat/backend/internal/service/websocket"
	"go.uber.org/zap"
)

type Bundle struct {
	Auth          *auth.AuthService
	Contacts      *contacts.ContactService
	Chats         *chats.ChatService
	Groups        *groups.GroupService
	Messages      *messages.MessageService
	Uploads       *uploads.UploadService
	Profile       *profile.ProfileService
	Notifications *notifications.NotificationService
	WSHub         *websocket.Hub
}

func NewBundle(cfg config.Config, repos repository.Bundle, logger *zap.Logger) Bundle {
	hub := websocket.NewHub(logger)
	
	// Start hub in background
	go hub.Run()

	return Bundle{
		Auth:          auth.NewAuthService(cfg, repos.Sessions, repos.Users, logger),
		Contacts:      contacts.NewContactService(repos.Contacts, repos.Users, logger),
		Chats:         chats.NewChatService(repos.Chats, repos.Messages, hub, logger),
		Groups:        groups.NewGroupService(repos.Chats, repos.Chats, hub, logger),
		Messages:      messages.NewMessageService(repos.Messages, repos.Attachments, hub, logger),
		Uploads:       uploads.NewUploadService(cfg.UploadDir, repos.Attachments, logger),
		Profile:       profile.NewProfileService(repos.Users, logger),
		Notifications: notifications.NewNotificationService(repos.Notifications, logger),
		WSHub:         hub,
	}
}
