package attachment

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachmentRepo struct {
	db *pgxpool.Pool
}

func NewAttachmentRepo(db *pgxpool.Pool) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

func (r *AttachmentRepo) Create(ctx context.Context, upload *model.Upload) error {
	query := `
		INSERT INTO uploads (id, user_id, file_name, file_size, content_type, storage_path, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, upload.ID, upload.UserID, upload.FileName, upload.FileSize, upload.ContentType, upload.StoragePath, upload.Status)
	return err
}

func (r *AttachmentRepo) GetByID(ctx context.Context, id string) (*model.Upload, error) {
	query := `
		SELECT id, user_id, file_name, file_size, content_type, storage_path, status, created_at
		FROM uploads WHERE id = $1
	`
	var upload model.Upload
	err := r.db.QueryRow(ctx, query, id).Scan(
		&upload.ID, &upload.UserID, &upload.FileName, &upload.FileSize,
		&upload.ContentType, &upload.StoragePath, &upload.Status, &upload.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *AttachmentRepo) CreateAttachment(ctx context.Context, attachment *model.Attachment) error {
	query := `
		INSERT INTO attachments (id, message_id, upload_id, file_name, file_size, content_type, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, attachment.ID, attachment.MessageID, attachment.UploadID, attachment.FileName, attachment.FileSize, attachment.ContentType, attachment.URL)
	return err
}

func (r *AttachmentRepo) GetByMessageID(ctx context.Context, messageID string) ([]model.Attachment, error) {
	query := `
		SELECT id, message_id, upload_id, file_name, file_size, content_type, url, created_at
		FROM attachments WHERE message_id = $1
	`
	rows, err := r.db.Query(ctx, query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []model.Attachment
	for rows.Next() {
		var a model.Attachment
		if err := rows.Scan(&a.ID, &a.MessageID, &a.UploadID, &a.FileName, &a.FileSize, &a.ContentType, &a.URL, &a.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}
