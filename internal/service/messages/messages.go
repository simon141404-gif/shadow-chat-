package messages

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/message"
	"github.com/yourorg/shadowchat/backend/internal/repository/attachment"
	"github.com/yourorg/shadowchat/backend/internal/service/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageService struct {
	messages    *message.MessageRepo
	attachments *attachment.AttachmentRepo
	hub         *websocket.Hub
	logger      *zap.Logger
}

func NewMessageService(messages *message.MessageRepo, attachments *attachment.AttachmentRepo, hub *websocket.Hub, logger *zap.Logger) *MessageService {
	return &MessageService{messages: messages, attachments: attachments, hub: hub, logger: logger}
}

type SendMessageRequest struct {
	ClientMsgID      string             `json:"clientMsgId"`
	MessageType      model.MessageType  `json:"messageType"`
	Content         string             `json:"content"`
	ReplyToMessageID *string           `json:"replyToMessageId,omitempty"`
	Attachments     []string           `json:"attachments,omitempty"`
}

func (s *MessageService) List(ctx context.Context, chatID string, limit, offset int) ([]model.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.messages.ListByChatID(ctx, chatID, limit, offset)
}

func (s *MessageService) Send(ctx context.Context, chatID, senderUserID string, req SendMessageRequest) (*model.Message, error) {
	msg := &model.Message{
		ID:               uuid.New().String(),
		ChatID:           chatID,
		SenderUserID:     senderUserID,
		ClientMsgID:      req.ClientMsgID,
		MessageType:      req.MessageType,
		Content:          req.Content,
		ReplyToMessageID: req.ReplyToMessageID,
		CreatedAt:        time.Now(),
	}

	if err := s.messages.Create(ctx, msg); err != nil {
		return nil, err
	}

	// Broadcast to WebSocket clients
	s.hub.Broadcast(chatID, senderUserID, msg)

	return msg, nil
}

func (s *MessageService) Edit(ctx context.Context, messageID, userID, content string) (*model.Message, error) {
	msg, err := s.messages.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if msg.SenderUserID != userID {
		return nil, ErrUnauthorized
	}

	msg.Content = content
	now := time.Now()
	msg.EditedAt = &now

	if err := s.messages.Update(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) Delete(ctx context.Context, messageID, userID string) error {
	msg, err := s.messages.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	if msg.SenderUserID != userID {
		return ErrUnauthorized
	}

	return s.messages.Delete(ctx, messageID)
}

func (s *MessageService) GetAttachments(ctx context.Context, messageID string) ([]model.Attachment, error) {
	return s.attachments.GetByMessageID(ctx, messageID)
}

var ErrUnauthorized = &AppError{Message: "unauthorized"}

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
