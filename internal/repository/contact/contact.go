package contact

import (
	"context"

	"github.com/yourorg/shadowchat/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepo struct {
	db *pgxpool.Pool
}

func NewContactRepo(db *pgxpool.Pool) *ContactRepo {
	return &ContactRepo{db: db}
}

func (r *ContactRepo) Create(ctx context.Context, contact *model.Contact) error {
	query := `
		INSERT INTO contacts (id, owner_user_id, contact_user_id, display_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (owner_user_id, contact_user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, contact.ID, contact.OwnerUserID, contact.ContactUserID, contact.DisplayName)
	return err
}

func (r *ContactRepo) ListByUserID(ctx context.Context, userID string) ([]model.Contact, error) {
	query := `
		SELECT c.id, c.owner_user_id, c.contact_user_id, c.display_name, c.created_at
		FROM contacts c
		WHERE c.owner_user_id = $1
		ORDER BY c.display_name, c.created_at
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var c model.Contact
		if err := rows.Scan(&c.ID, &c.OwnerUserID, &c.ContactUserID, &c.DisplayName, &c.CreatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *ContactRepo) Get(ctx context.Context, ownerUserID, contactUserID string) (*model.Contact, error) {
	query := `
		SELECT id, owner_user_id, contact_user_id, display_name, created_at
		FROM contacts WHERE owner_user_id = $1 AND contact_user_id = $2
	`
	var contact model.Contact
	err := r.db.QueryRow(ctx, query, ownerUserID, contactUserID).Scan(
		&contact.ID, &contact.OwnerUserID, &contact.ContactUserID, &contact.DisplayName, &contact.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *ContactRepo) Delete(ctx context.Context, ownerUserID, contactUserID string) error {
	query := `DELETE FROM contacts WHERE owner_user_id = $1 AND contact_user_id = $2`
	_, err := r.db.Exec(ctx, query, ownerUserID, contactUserID)
	return err
}
