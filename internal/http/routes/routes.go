package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/http/handlers"
	"github.com/yourorg/shadowchat/backend/internal/http/middleware"
	"github.com/yourorg/shadowchat/backend/internal/service"
	"go.uber.org/zap"
)

func NewRouter(cfg config.Config, svc service.Bundle, logger *zap.Logger) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestID(), middleware.Logging(logger), cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Request-Id"},
		AllowCredentials: true,
	}))

	h := handlers.NewBundle(svc, cfg, logger)

	v1 := r.Group("/v1")
	{
		v1.POST("/auth/anonymous", h.Auth.Anonymous)
		v1.POST("/auth/refresh", h.Auth.Refresh)

		secured := v1.Group("")
		secured.Use(middleware.Auth(cfg.JWTSecret))
		{
			// Chats
			secured.GET("/chats", h.Chats.List)
			secured.POST("/chats", h.Chats.Create)
			secured.GET("/chats/:chatId", h.Chats.Get)
			secured.GET("/chats/:chatId/messages", h.Chats.ListMessages)
			secured.POST("/chats/:chatId/messages", h.Chats.SendMessage)

			// Messages
			secured.PATCH("/messages/:messageId", h.Messages.Edit)
			secured.DELETE("/messages/:messageId", h.Messages.Delete)

			// Contacts
			secured.GET("/contacts", h.Contacts.List)
			secured.POST("/contacts/share", h.Contacts.Share)
			secured.POST("/contacts/qr", h.Contacts.QR)

			// Groups
			secured.POST("/groups", h.Groups.Create)
			secured.POST("/groups/:chatId/members", h.Groups.AddMembers)

			// Uploads
			secured.POST("/uploads", h.Uploads.Create)
			secured.GET("/uploads/:uploadId", h.Uploads.Get)

			// Profile
			secured.GET("/profile", h.Profile.Get)
			secured.PATCH("/profile", h.Profile.Update)

			// Notifications
			secured.GET("/notifications", h.Notifications.List)
			secured.POST("/notifications/push-token", h.Notifications.RegisterPushToken)

			// WebSocket
			secured.GET("/ws", h.WS.Serve)
		}
	}

	// Serve uploaded files
	r.Static("/uploads", cfg.UploadDir)

	return r
}
