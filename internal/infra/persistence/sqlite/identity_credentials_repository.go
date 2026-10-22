package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

type IdentityCredentialsRepository struct{ db *sql.DB }

func NewIdentityCredentialsRepository(db *sql.DB) *IdentityCredentialsRepository {
	return &IdentityCredentialsRepository{db: db}
}

func (r *IdentityCredentialsRepository) Save(
	ctx context.Context,
	value identity.IdentityCredentials,
) error {
	validated, err := identity.NewIdentityCredentials(
		value.IdentityID,
		value.AccessTokenEncrypted,
		value.RefreshTokenEncrypted,
		value.AccessTokenExpiresAt,
		value.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		`INSERT INTO identity_credentials
		(identity_id, access_token_encrypted, refresh_token_encrypted, access_token_expires_at, updated_at)
		VALUES (?, ?, ?, ?, ?) ON CONFLICT(identity_id) DO UPDATE SET
		access_token_encrypted = excluded.access_token_encrypted,
		refresh_token_encrypted = excluded.refresh_token_encrypted,
		access_token_expires_at = excluded.access_token_expires_at, updated_at = excluded.updated_at`,
		validated.IdentityID,
		validated.AccessTokenEncrypted,
		validated.RefreshTokenEncrypted,
		validated.AccessTokenExpiresAt.UTC().Format(time.RFC3339Nano),
		validated.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (r *IdentityCredentialsRepository) FindByID(
	ctx context.Context,
	id string,
) (identity.IdentityCredentials, error) {
	var identityID string
	var accessTokenEncrypted string
	var refreshTokenEncrypted string
	var expiry string
	var updated string
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT identity_id, access_token_encrypted, refresh_token_encrypted,
		access_token_expires_at, updated_at FROM identity_credentials WHERE identity_id = ?`,
		id,
	).Scan(
		&identityID,
		&accessTokenEncrypted,
		&refreshTokenEncrypted,
		&expiry,
		&updated,
	)
	if err != nil {
		return identity.IdentityCredentials{}, err
	}

	expiresAt, err := time.Parse(
		time.RFC3339Nano,
		expiry,
	)
	if err != nil {
		return identity.IdentityCredentials{}, err
	}
	updatedAt, err := time.Parse(
		time.RFC3339Nano,
		updated,
	)
	if err != nil {
		return identity.IdentityCredentials{}, err
	}

	return identity.NewIdentityCredentials(
		identityID,
		accessTokenEncrypted,
		refreshTokenEncrypted,
		expiresAt,
		updatedAt,
	)
}

func (r *IdentityCredentialsRepository) FindLatest(ctx context.Context) (identity.IdentityCredentials, error) {
	var id string
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT identity_id FROM identity_credentials
		ORDER BY julianday(updated_at) DESC, updated_at DESC, identity_id ASC LIMIT 1`,
	).Scan(&id)
	if err != nil {
		return identity.IdentityCredentials{}, err
	}
	return r.FindByID(
		ctx,
		id,
	)
}

var _ identity.IdentityCredentialsRepository = (*IdentityCredentialsRepository)(nil)
