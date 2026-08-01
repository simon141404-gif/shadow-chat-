package handlers

import (
	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/service"
	authHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/auth"
	chatsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/chats"
	contactsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/contacts"
	groupsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/groups"
	messagesHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/messages"
	notificationsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/notifications"
	profileHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/profile"
	uploadsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/uploads"
	wsHandler "github.com/yourorg/shadowchat/backend/internal/http/handlers/ws"
	"go.uber.org/zap"
)

type Bundle struct {
	Auth          *authHandler.AuthHandler
	Chats         *chatsHandler.ChatsHandler
	Contacts      *contactsHandler.ContactsHandler
	Groups        *groupsHandler.GroupsHandler
	Messages      *messagesHandler.MessageHandler
	Uploads       *uploadsHandler.UploadsHandler
	Profile       *profileHandler.ProfileHandler
	Notifications *notificationsHandler.NotificationsHandler
	WS            *wsHandler.WSHandler
}

func NewBundle(svc service.Bundle, cfg config.Config, logger *zap.Logger) Bundle {
	return Bundle{
		Auth:          authHandler.NewAuthHandler(svc.Auth),
		Chats:         chatsHandler.NewChatsHandler(svc.Chats, svc.Messages),
		Contacts:      contactsHandler.NewContactsHandler(svc.Contacts),
		Groups:        groupsHandler.NewGroupsHandler(svc.Groups),
		Messages:      messagesHandler.NewMessageHandler(svc.Messages),
		Uploads:       uploadsHandler.NewUploadsHandler(svc.Uploads),
		Profile:       profileHandler.NewProfileHandler(svc.Profile),
		Notifications: notificationsHandler.NewNotificationsHandler(svc.Notifications),
		WS:            wsHandler.NewWSHandler(svc.WSHub),
	}
}
