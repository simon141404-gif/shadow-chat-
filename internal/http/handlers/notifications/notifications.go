package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/notifications"
)

type NotificationsHandler struct {
	svc *notifications.NotificationService
}

func NewNotificationsHandler(svc *notifications.NotificationService) *NotificationsHandler {
	return &NotificationsHandler{svc: svc}
}

type RegisterPushTokenRequest struct {
	Token string `json:"token"`
}

func (h *NotificationsHandler) RegisterPushToken(c *gin.Context) {
	userID := c.GetString("userId")
	var req RegisterPushTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.RegisterPushToken(c.Request.Context(), userID, notifications.RegisterPushTokenRequest{
		Token: req.Token,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register push token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "push token registered"})
}

func (h *NotificationsHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	notifications, err := h.svc.List(c.Request.Context(), userID, 20, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}
