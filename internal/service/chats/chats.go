package chats

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/chat"
	"github.com/yourorg/shadowchat/backend/internal/service/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChatService struct {
	chats     *chat.ChatRepo
	messages  interface{ ListByChatID(ctx context.Context, chatID string, limit, offset int) ([]model.Message, error) }
	hub       *websocket.Hub
	logger    *zap.Logger
}

func NewChatService(chats *chat.ChatRepo, messages interface{ ListByChatID(ctx context.Context, chatID string, limit, offset int) ([]model.Message, error) }, hub *websocket.Hub, logger *zap.Logger) *ChatService {
	return &ChatService{chats: chats, messages: messages, hub: hub, logger: logger}
}

type CreateChatRequest struct {
	Type    model.ChatType `json:"type"`
	Name    string        `json:"name,omitempty"`
	Members []string      `json:"members,omitempty"`
}

func (s *ChatService) List(ctx context.Context, userID string) ([]model.Chat, error) {
	return s.chats.ListByUserID(ctx, userID)
}

func (s *ChatService) Create(ctx context.Context, creatorUserID string, req CreateChatRequest) (*model.Chat, error) {
	chat := &model.Chat{
		ID:        uuid.New().String(),
		Type:      req.Type,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chats.Create(ctx, chat); err != nil {
		return nil, err
	}

	// Add creator as member
	if err := s.chats.AddMember(ctx, &model.ChatMember{
		ChatID:   chat.ID,
		UserID:   creatorUserID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	// Add other members
	for _, memberID := range req.Members {
		s.chats.AddMember(ctx, &model.ChatMember{
			ChatID:    chat.ID,
			UserID:    memberID,
			Role:      "member",
			JoinedAt:  time.Now(),
			InvitedBy: creatorUserID,
		})
	}

	return chat, nil
}

func (s *ChatService) Get(ctx context.Context, chatID, userID string) (*model.Chat, error) {
	// Verify membership
	isMember, err := s.chats.IsMember(ctx, chatID, userID)
	if err != nil || !isMember {
		return nil, err
	}

	return s.chats.GetByID(ctx, chatID)
}

func (s *ChatService) GetMembers(ctx context.Context, chatID string) ([]model.ChatMember, error) {
	return s.chats.GetMembers(ctx, chatID)
}
