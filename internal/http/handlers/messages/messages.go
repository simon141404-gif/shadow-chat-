package messages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/messages"
)

type MessageHandler struct {
	svc *messages.MessageService
}

func NewMessageHandler(svc *messages.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

type EditRequest struct {
	Content string `json:"content"`
}

func (h *MessageHandler) Edit(c *gin.Context) {
	userID := c.GetString("userId")
	messageID := c.Param("messageId")
	var req EditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	msg, err := h.svc.Edit(c.Request.Context(), messageID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to edit message"})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func (h *MessageHandler) Delete(c *gin.Context) {
	userID := c.GetString("userId")
	messageID := c.Param("messageId")

	if err := h.svc.Delete(c.Request.Context(), messageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
