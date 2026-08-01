package contacts

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/contact"
	"github.com/yourorg/shadowchat/backend/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContactService struct {
	contacts *contact.ContactRepo
	users    *user.UserRepo
	logger   *zap.Logger
}

func NewContactService(contacts *contact.ContactRepo, users *user.UserRepo, logger *zap.Logger) *ContactService {
	return &ContactService{contacts: contacts, users: users, logger: logger}
}

type ShareContactRequest struct {
	PublicID string `json:"publicId"`
}

func (s *ContactService) Share(ctx context.Context, ownerUserID string, req ShareContactRequest) (*model.Contact, error) {
	// Find user by public ID
	user, err := s.users.GetByPublicID(ctx, req.PublicID)
	if err != nil {
		return nil, err
	}

	contact := &model.Contact{
		ID:            uuid.New().String(),
		OwnerUserID:   ownerUserID,
		ContactUserID: user.ID,
		DisplayName:   user.DisplayName,
	}

	if err := s.contacts.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactService) List(ctx context.Context, userID string) ([]model.Contact, error) {
	return s.contacts.ListByUserID(ctx, userID)
}

func (s *ContactService) GetQRCode(ctx context.Context, userID string) (string, error) {
	// Generate QR code data for user
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	return user.PublicID, nil
}
