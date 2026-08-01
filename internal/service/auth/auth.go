package auth

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/yourorg/shadowchat/backend/internal/repository/session"
	"github.com/yourorg/shadowchat/backend/internal/repository/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg      config.Config
	sessions *session.SessionRepo
	users    *user.UserRepo
	logger   *zap.Logger
}

func NewAuthService(cfg config.Config, sessions *session.SessionRepo, users *user.UserRepo, logger *zap.Logger) *AuthService {
	return &AuthService{
		cfg:      cfg,
		sessions: sessions,
		users:    users,
		logger:   logger,
	}
}

type Claims struct {
	UserID string `json:"userId"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func (s *AuthService) CreateAnonymousIdentity(ctx context.Context) (userID, publicID string, err error) {
	publicID = generatePublicID()
	user := &model.User{
		ID:       uuid.New().String(),
		PublicID: publicID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return "", "", err
	}

	// Create default session
	sessionID := uuid.New().String()
	jti := uuid.New().String()
	refreshTokenHash, _ := hashToken(jti)

	session := &model.Session{
		ID:               sessionID,
		UserID:           user.ID,
		JTI:              jti,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.sessions.Create(ctx, session); err != nil {
		return "", "", err
	}

	return user.ID, publicID, nil
}

func (s *AuthService) GenerateJWT(userID, jti string) (string, error) {
	claims := Claims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token and create new session
	jti := uuid.New().String()
	userID, err := s.extractUserIDFromRefresh(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Create new session
	sessionID := uuid.New().String()
	refreshTokenHash, _ := hashToken(jti)

	session := &model.Session{
		ID:               sessionID,
		UserID:           userID,
		JTI:              jti,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.sessions.Create(ctx, session); err != nil {
		return "", "", err
	}

	// Generate new JWT
	jwtToken, err := s.GenerateJWT(userID, jti)
	if err != nil {
		return "", "", err
	}

	return jwtToken, jti, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, jti string) error {
	if err := s.sessions.Revoke(ctx, jti); err != nil {
		return err
	}
	return s.sessions.StoreTokenBlacklist(ctx, jti, 24*time.Hour)
}

func (s *AuthService) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	return s.sessions.IsTokenBlacklisted(ctx, jti)
}

func (s *AuthService) extractUserIDFromRefresh(refreshToken string) (string, error) {
	// In production, validate the refresh token properly
	// For now, return a placeholder - implement proper validation
	return "", nil
}

func generatePublicID() string {
	return uuid.New().String()[:8]
}

func hashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(hash), err
}
