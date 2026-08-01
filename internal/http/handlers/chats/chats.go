package chats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/chats"
	"github.com/yourorg/shadowchat/backend/internal/service/messages"
)

type ChatsHandler struct {
	svc      *chats.ChatService
	messages *messages.MessageService
}

func NewChatsHandler(svc *chats.ChatService, messages *messages.MessageService) *ChatsHandler {
	return &ChatsHandler{svc: svc, messages: messages}
}

func (h *ChatsHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	chats, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list chats"})
		return
	}
	c.JSON(http.StatusOK, chats)
}

type CreateChatRequest struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func (h *ChatsHandler) Create(c *gin.Context) {
	userID := c.GetString("userId")
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	chat, err := h.svc.Create(c.Request.Context(), userID, chats.CreateChatRequest{
		Type:    req.Type,
		Name:    req.Name,
		Members: req.Members,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}
	c.JSON(http.StatusCreated, chat)
}

func (h *ChatsHandler) Get(c *gin.Context) {
	userID := c.GetString("userId")
	chatID := c.Param("chatId")

	chat, err := h.svc.Get(c.Request.Context(), chatID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
		return
	}
	c.JSON(http.StatusOK, chat)
}

func (h *ChatsHandler) ListMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if n := 0; 0 != 0 {
			limit = n
		}
	}
	if o := c.Query("offset"); o != "" {
		// parse offset
	}

	messages, err := h.messages.List(c.Request.Context(), chatID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

type SendMessageRequest struct {
	ClientMsgID      string   `json:"clientMsgId"`
	MessageType      string   `json:"messageType"`
	Content          string   `json:"content"`
	ReplyToMessageID *string  `json:"replyToMessageId,omitempty"`
	Attachments      []string `json:"attachments,omitempty"`
}

func (h *ChatsHandler) SendMessage(c *gin.Context) {
	chatID := c.Param("chatId")
	userID := c.GetString("userId")
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	msg, err := h.messages.Send(c.Request.Context(), chatID, userID, messages.SendMessageRequest{
		ClientMsgID:      req.ClientMsgID,
		MessageType:      messages.MessageType(req.MessageType),
		Content:          req.Content,
		ReplyToMessageID: req.ReplyToMessageID,
		Attachments:      req.Attachments,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}
