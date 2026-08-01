package model

import "time"

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeVideo    MessageType = "video"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeDocument MessageType = "document"
	MessageTypeVoice    MessageType = "voice"
)

type Message struct {
	ID               string       `json:"id"`
	ChatID           string       `json:"chatId"`
	SenderUserID     string       `json:"senderUserId"`
	ClientMsgID      string       `json:"clientMsgId"`
	Ciphertext       []byte       `json:"-"`
	MessageType      MessageType  `json:"messageType"`
	ReplyToMessageID *string     `json:"replyToMessageId,omitempty"`
	Content          string       `json:"content"`
	CreatedAt        time.Time    `json:"createdAt"`
	EditedAt         *time.Time  `json:"editedAt,omitempty"`
	DeletedAt        *time.Time  `json:"deletedAt,omitempty"`
	ExpiresAt        *time.Time  `json:"expiresAt,omitempty"`
}
