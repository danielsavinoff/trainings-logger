package sqlite

import (
	"context"
	"database/sql"

	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

type IdentityRepository struct {
	db *sql.DB
}

func NewIdentityRepository(db *sql.DB) *IdentityRepository {
	return &IdentityRepository{db: db}
}

func (r *IdentityRepository) Save(
	ctx context.Context,
	value identity.Identity,
) error {
	validated, err := identity.NewIdentity(
		value.ID,
		value.Provider,
	)
	if err != nil {
		return err
	}

	_, err = persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		`INSERT INTO identities (id, provider) VALUES (?, ?)
		ON CONFLICT(id) DO UPDATE SET provider = excluded.provider`,
		validated.ID,
		validated.Provider,
	)
	return err
}

func (r *IdentityRepository) FindByID(
	ctx context.Context,
	id string,
) (identity.Identity, error) {
	var storedID string
	var provider string
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		"SELECT id, provider FROM identities WHERE id = ?",
		id,
	).Scan(
		&storedID,
		&provider,
	)
	if err != nil {
		return identity.Identity{}, err
	}

	validatedProvider, err := identity.NewProvider(provider)
	if err != nil {
		return identity.Identity{}, err
	}

	return identity.NewIdentity(
		storedID,
		validatedProvider,
	)
}

var _ identity.IdentityRepository = (*IdentityRepository)(nil)
