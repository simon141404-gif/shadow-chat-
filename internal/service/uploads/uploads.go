package uploads

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/attachment"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UploadService struct {
	uploadDir  string
	attachments *attachment.AttachmentRepo
	logger     *zap.Logger
}

func NewUploadService(uploadDir string, attachments *attachment.AttachmentRepo, logger *zap.Logger) *UploadService {
	return &UploadService{uploadDir: uploadDir, attachments: attachments, logger: logger}
}

type CreateUploadRequest struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	ContentType string `json:"contentType"`
}

func (s *UploadService) Create(ctx context.Context, userID string, req CreateUploadRequest) (*model.Upload, error) {
	upload := &model.Upload{
		ID:          uuid.New().String(),
		UserID:      userID,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
		StoragePath: filepath.Join(userID, upload.ID),
		Status:      "pending",
	}

	if err := s.attachments.Create(ctx, upload); err != nil {
		return nil, err
	}

	// Ensure upload directory exists
	fullPath := filepath.Join(s.uploadDir, upload.StoragePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}

	return upload, nil
}

func (s *UploadService) SaveFile(ctx context.Context, uploadID string, reader io.Reader) error {
	upload, err := s.attachments.GetByID(ctx, uploadID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(s.uploadDir, upload.StoragePath)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	// Update status to completed
	upload.Status = "completed"
	return nil
}

func (s *UploadService) GetUploadURL(upload *model.Upload) string {
	return fmt.Sprintf("/uploads/%s", upload.StoragePath)
}

func (s *UploadService) GetFilePath(upload *model.Upload) string {
	return filepath.Join(s.uploadDir, upload.StoragePath)
}

func (s *UploadService) GetUpload(ctx context.Context, uploadID string) (*model.Upload, error) {
	return s.attachments.GetByID(ctx, uploadID)
}
