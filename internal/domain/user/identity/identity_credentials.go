package identity

import (
	"strings"
	"time"
)

type IdentityCredentials struct {
	IdentityID            string
	AccessTokenEncrypted  string
	RefreshTokenEncrypted string
	AccessTokenExpiresAt  time.Time
	UpdatedAt             time.Time
}

func NewIdentityCredentials(
	identityID string,
	accessTokenEncrypted string,
	refreshTokenEncrypted string,
	accessTokenExpiresAt time.Time,
	updatedAt time.Time,
) (IdentityCredentials, error) {
	if strings.TrimSpace(identityID) == "" ||
		strings.TrimSpace(accessTokenEncrypted) == "" ||
		strings.TrimSpace(refreshTokenEncrypted) == "" ||
		accessTokenExpiresAt.IsZero() ||
		updatedAt.IsZero() {
		return IdentityCredentials{}, ErrInvalidCredentials
	}

	return IdentityCredentials{
		IdentityID:            identityID,
		AccessTokenEncrypted:  accessTokenEncrypted,
		RefreshTokenEncrypted: refreshTokenEncrypted,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		UpdatedAt:             updatedAt,
	}, nil
}
