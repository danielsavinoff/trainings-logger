package identity

import (
	"strings"
	"time"
)

type IdentityCredentials struct {
	IdentityID           string    `json:"identity_id"`
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func NewIdentityCredentials(
	identityID string,
	accessToken string,
	refreshToken string,
	accessTokenExpiresAt time.Time,
	updatedAt time.Time,
) (IdentityCredentials, error) {
	if strings.TrimSpace(identityID) == "" ||
		strings.TrimSpace(accessToken) == "" ||
		strings.TrimSpace(refreshToken) == "" ||
		accessTokenExpiresAt.IsZero() ||
		updatedAt.IsZero() {
		return IdentityCredentials{}, ErrInvalidCredentials
	}

	return IdentityCredentials{
		IdentityID:           identityID,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		UpdatedAt:            updatedAt,
	}, nil
}
