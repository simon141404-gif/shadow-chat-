package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/groups"
)

type GroupsHandler struct {
	svc *groups.GroupService
}

func NewGroupsHandler(svc *groups.GroupService) *GroupsHandler {
	return &GroupsHandler{svc: svc}
}

type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func (h *GroupsHandler) Create(c *gin.Context) {
	userID := c.GetString("userId")
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	group, err := h.svc.Create(c.Request.Context(), userID, groups.CreateGroupRequest{
		Name:    req.Name,
		Members: req.Members,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}
	c.JSON(http.StatusCreated, group)
}

type AddMembersRequest struct {
	Members []string `json:"members"`
}

func (h *GroupsHandler) AddMembers(c *gin.Context) {
	userID := c.GetString("userId")
	chatID := c.Param("chatId")
	var req AddMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.AddMembers(c.Request.Context(), chatID, userID, req.Members); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add members"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "members added"})
}
