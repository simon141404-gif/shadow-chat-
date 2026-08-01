package user

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, public_id, display_name, avatar_url, bio)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.PublicID, user.DisplayName, user.AvatarURL, user.Bio)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, public_id, display_name, avatar_url, bio, created_at, updated_at, deleted_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.PublicID, &user.DisplayName, &user.AvatarURL, &user.Bio,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByPublicID(ctx context.Context, publicID string) (*model.User, error) {
	query := `
		SELECT id, public_id, display_name, avatar_url, bio, created_at, updated_at, deleted_at
		FROM users WHERE public_id = $1 AND deleted_at IS NULL
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, publicID).Scan(
		&user.ID, &user.PublicID, &user.DisplayName, &user.AvatarURL, &user.Bio,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET display_name = $1, avatar_url = $2, bio = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, user.DisplayName, user.AvatarURL, user.Bio, user.ID)
	return err
}
