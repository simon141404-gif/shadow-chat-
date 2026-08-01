package model

import "time"

type Session struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	JTI             string     `json:"jti"`
	RefreshTokenHash string   `json:"-"`
	ExpiresAt       time.Time  `json:"expiresAt"`
	RevokedAt       *time.Time `json:"revokedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}
