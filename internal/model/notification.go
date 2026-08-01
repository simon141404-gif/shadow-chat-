package model

import "time"

type NotificationType string

const (
	NotificationTypeMessage NotificationType = "message"
	NotificationTypeContact NotificationType = "contact"
	NotificationTypeGroup   NotificationType = "group"
)

type Notification struct {
	ID        string           `json:"id"`
	UserID    string           `json:"userId"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	Data      string           `json:"data,omitempty"`
	ReadAt    *time.Time      `json:"readAt,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
}

type PushToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}
