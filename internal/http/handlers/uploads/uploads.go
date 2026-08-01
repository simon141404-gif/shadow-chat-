package uploads

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/uploads"
)

type UploadsHandler struct {
	svc *uploads.UploadService
}

func NewUploadsHandler(svc *uploads.UploadService) *UploadsHandler {
	return &UploadsHandler{svc: svc}
}

type CreateUploadRequest struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	ContentType string `json:"contentType"`
}

func (h *UploadsHandler) Create(c *gin.Context) {
	userID := c.GetString("userId")
	var req CreateUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	upload, err := h.svc.Create(c.Request.Context(), userID, uploads.CreateUploadRequest{
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":           upload.ID,
		"uploadUrl":    "/v1/uploads/" + upload.ID + "/data",
		"storagePath":  upload.StoragePath,
	})
}

func (h *UploadsHandler) Get(c *gin.Context) {
	uploadID := c.Param("uploadId")
	upload, err := h.svc.GetUpload(c.Request.Context(), uploadID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "upload not found"})
		return
	}
	c.JSON(http.StatusOK, upload)
}
