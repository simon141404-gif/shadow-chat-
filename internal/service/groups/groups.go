package groups

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/chat"
	"github.com/yourorg/shadowchat/backend/internal/service/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GroupService struct {
	groups *chat.ChatRepo
	chats  *chat.ChatRepo
	hub    *websocket.Hub
	logger *zap.Logger
}

func NewGroupService(groups *chat.ChatRepo, chats *chat.ChatRepo, hub *websocket.Hub, logger *zap.Logger) *GroupService {
	return &GroupService{groups: groups, chats: chats, hub: hub, logger: logger}
}

type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func (s *GroupService) Create(ctx context.Context, creatorUserID string, req CreateGroupRequest) (*model.Chat, error) {
	chat := &model.Chat{
		ID:        uuid.New().String(),
		Type:      model.ChatTypeGroup,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.groups.Create(ctx, chat); err != nil {
		return nil, err
	}

	// Add creator as owner
	if err := s.groups.AddMember(ctx, &model.ChatMember{
		ChatID:   chat.ID,
		UserID:   creatorUserID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	// Add members
	for _, memberID := range req.Members {
		s.groups.AddMember(ctx, &model.ChatMember{
			ChatID:    chat.ID,
			UserID:    memberID,
			Role:      "member",
			JoinedAt:  time.Now(),
			InvitedBy: creatorUserID,
		})
	}

	// Broadcast group creation
	s.hub.Broadcast(chat.ID, creatorUserID, map[string]string{
		"type": "group_created",
		"name": req.Name,
	})

	return chat, nil
}

func (s *GroupService) AddMembers(ctx context.Context, chatID, inviterUserID string, memberIDs []string) error {
	for _, memberID := range memberIDs {
		if err := s.groups.AddMember(ctx, &model.ChatMember{
			ChatID:    chatID,
			UserID:    memberID,
			Role:      "member",
			JoinedAt:  time.Now(),
			InvitedBy: inviterUserID,
		}); err != nil {
			return err
		}
	}

	// Broadcast member addition
	s.hub.Broadcast(chatID, inviterUserID, map[string]interface{}{
		"type":    "members_added",
		"members": memberIDs,
	})

	return nil
}

func (s *GroupService) RemoveMember(ctx context.Context, chatID, removerUserID, memberID string) error {
	// In production, verify permissions and perform actual removal
	s.hub.Broadcast(chatID, removerUserID, map[string]interface{}{
		"type":      "member_removed",
		"member_id": memberID,
	})
	return nil
}
