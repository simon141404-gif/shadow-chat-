package session

import (
	"context"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SessionRepo struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewSessionRepo(db *pgxpool.Pool, redis *redis.Client) *SessionRepo {
	return &SessionRepo{db: db, redis: redis}
}

func (r *SessionRepo) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, jti, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, session.ID, session.UserID, session.JTI, session.RefreshTokenHash, session.ExpiresAt)
	return err
}

func (r *SessionRepo) GetByJTI(ctx context.Context, jti string) (*model.Session, error) {
	query := `
		SELECT id, user_id, jti, refresh_token_hash, expires_at, revoked_at, created_at
		FROM sessions WHERE jti = $1
	`
	var session model.Session
	err := r.db.QueryRow(ctx, query, jti).Scan(
		&session.ID, &session.UserID, &session.JTI, &session.RefreshTokenHash,
		&session.ExpiresAt, &session.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepo) Revoke(ctx context.Context, jti string) error {
	query := `UPDATE sessions SET revoked_at = NOW() WHERE jti = $1`
	_, err := r.db.Exec(ctx, query, jti)
	return err
}

func (r *SessionRepo) StoreTokenBlacklist(ctx context.Context, jti string, exp time.Duration) error {
	return r.redis.Set(ctx, "blacklist:"+jti, "1", exp).Err()
}

func (r *SessionRepo) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	result, err := r.redis.Exists(ctx, "blacklist:"+jti).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
