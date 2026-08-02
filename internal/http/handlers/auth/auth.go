package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/auth"
)

type AuthHandler struct {
	svc *auth.AuthService
}

func NewAuthHandler(svc *auth.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type AnonymousResponse struct {
	UserID  string `json:"userId"`
	PublicID string `json:"publicId"`
	Token   string `json:"token"`
}

func (h *AuthHandler) Anonymous(c *gin.Context) {
	userID, publicID, err := h.svc.CreateAnonymousIdentity(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create identity: " + err.Error()})
		return
	}

	// Generate JWT for the new user
	token, err := h.svc.GenerateJWT(userID, publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AnonymousResponse{
		UserID:   userID,
		PublicID: publicID,
		Token:    token,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, jti, err := h.svc.RefreshSession(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"jti":   jti,
	})
}

type LogoutRequest struct {
	JTI string `json:"jti"`
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.RevokeSession(c.Request.Context(), req.JTI); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
