package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/contacts"
)

type ContactsHandler struct {
	svc *contacts.ContactService
}

func NewContactsHandler(svc *contacts.ContactService) *ContactsHandler {
	return &ContactsHandler{svc: svc}
}

type ShareRequest struct {
	PublicID string `json:"publicId"`
}

func (h *ContactsHandler) Share(c *gin.Context) {
	userID := c.GetString("userId")
	var req ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	contact, err := h.svc.Share(c.Request.Context(), userID, contacts.ShareContactRequest{
		PublicID: req.PublicID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to share contact"})
		return
	}
	c.JSON(http.StatusCreated, contact)
}

func (h *ContactsHandler) List(c *gin.Context) {
	userID := c.GetString("userId")
	contacts, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list contacts"})
		return
	}
	c.JSON(http.StatusOK, contacts)
}

type QRCodeResponse struct {
	QRCode string `json:"qrCode"`
}

func (h *ContactsHandler) QR(c *gin.Context) {
	userID := c.GetString("userId")
	qrCode, err := h.svc.GetQRCode(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR code"})
		return
	}
	c.JSON(http.StatusOK, QRCodeResponse{QRCode: qrCode})
}
