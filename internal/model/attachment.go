package model

import "time"

type Attachment struct {
	ID          string     `json:"id"`
	MessageID   string     `json:"messageId"`
	UploadID    string     `json:"uploadId"`
	FileName    string     `json:"fileName"`
	FileSize    int64      `json:"fileSize"`
	ContentType string     `json:"contentType"`
	URL         string     `json:"url"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type Upload struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	FileName    string     `json:"fileName"`
	FileSize    int64      `json:"fileSize"`
	ContentType string     `json:"contentType"`
	StoragePath string     `json:"storagePath"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
}
