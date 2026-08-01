package model

import "time"

type Contact struct {
	ID          string     `json:"id"`
	OwnerUserID string     `json:"ownerUserId"`
	ContactUserID string  `json:"contactUserId"`
	DisplayName string     `json:"displayName,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
