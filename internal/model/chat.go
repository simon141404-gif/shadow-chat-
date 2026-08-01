package model

import "time"

type ChatType string

const (
	ChatTypeDirect ChatType = "direct"
	ChatTypeGroup ChatType = "group"
)

type Chat struct {
	ID          string     `json:"id"`
	Type        ChatType  `json:"type"`
	Name        string    `json:"name,omitempty"`
	AvatarURL   string    `json:"avatarUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type ChatMember struct {
	ChatID     string    `json:"chatId"`
	UserID     string    `json:"userId"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joinedAt"`
	InvitedBy  string    `json:"invitedBy,omitempty"`
}
