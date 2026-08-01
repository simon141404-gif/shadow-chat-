package profile

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/user"
	"go.uber.org/zap"
)

type ProfileService struct {
	users  *user.UserRepo
	logger *zap.Logger
}

func NewProfileService(users *user.UserRepo, logger *zap.Logger) *ProfileService {
	return &ProfileService{users: users, logger: logger}
}

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Bio         string `json:"bio"`
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*model.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.users.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
