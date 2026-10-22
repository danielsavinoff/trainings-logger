package identity

import (
	"strings"
	"time"
)

type IdentityCredentials struct {
	IdentityID            string    `json:"identity_id"`
	AccessTokenEncrypted  string    `json:"access_token_encrypted"`
	RefreshTokenEncrypted string    `json:"refresh_token_encrypted"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	UpdatedAt             time.Time `json:"updated_at"`
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
