package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/profile"
)

type ProfileHandler struct {
	svc *profile.ProfileService
}

func NewProfileHandler(svc *profile.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) Get(c *gin.Context) {
	userID := c.GetString("userId")
	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Bio         string `json:"bio"`
}

func (h *ProfileHandler) Update(c *gin.Context) {
	userID := c.GetString("userId")
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, profile.UpdateProfileRequest{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		Bio:         req.Bio,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, user)
}
