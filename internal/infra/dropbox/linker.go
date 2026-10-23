// START_MODULE_CONTRACT
//   PURPOSE: Authorize a real Dropbox account with desktop PKCE.
//   SCOPE: Authorization URLs, browser launch, code exchange, and local credentials.
//   DEPENDS: M-DOMAIN-IDENTITY, M-APP-QUESTIONNAIRE
//   LINKS: M-INFRA-DROPBOX-AUTH, V-M-INFRA-DROPBOX-AUTH
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   Linker - hold Dropbox configuration and external adapters
//   NewLinker - configure real Dropbox endpoints and browser launcher
//   Begin - create a fresh PKCE authorization session
//   exchange - exchange an authorization code for account credentials
// END_MODULE_MAP

package dropbox

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
)

// START_BLOCK_DROPBOX_AUTHORIZATION
type Linker struct {
	appKey      string
	tokenURL    string
	client      *http.Client
	openBrowser func(context.Context, string) error
}

func NewLinker(appKey string) *Linker {
	return &Linker{
		appKey:   strings.TrimSpace(appKey),
		tokenURL: "https://api.dropboxapi.com/oauth2/token",
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(
				_ *http.Request,
				_ []*http.Request,
			) error {
				return http.ErrUseLastResponse
			},
		},
		openBrowser: openBrowser,
	}
}

func (l *Linker) Begin(
	ctx context.Context,
	provider identity.Provider,
) (ports.IdentityAuthorization, error) {
	if err := ctx.Err(); err != nil {
		return ports.IdentityAuthorization{}, err
	}
	if provider != identity.ProviderDropbox {
		return ports.IdentityAuthorization{}, identity.ErrInvalidProvider
	}
	if l.appKey == "" {
		return ports.IdentityAuthorization{}, errors.New("Dropbox sign-in is not configured")
	}

	random := make(
		[]byte,
		32,
	)
	if _, err := rand.Read(random); err != nil {
		return ports.IdentityAuthorization{}, err
	}
	verifier := base64.RawURLEncoding.EncodeToString(random)
	challenge := sha256.Sum256([]byte(verifier))
	parameters := url.Values{
		"client_id":             {l.appKey},
		"response_type":         {"code"},
		"token_access_type":     {"offline"},
		"code_challenge_method": {"S256"},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(challenge[:])},
		"scope":                 {"account_info.read files.content.read files.content.write"},
		"force_reapprove":       {"true"},
	}
	authorizationURL := "https://www.dropbox.com/oauth2/authorize?" + parameters.Encode()
	opened := l.openBrowser(
		ctx,
		authorizationURL,
	) == nil
	if err := ctx.Err(); err != nil {
		return ports.IdentityAuthorization{}, err
	}

	return ports.IdentityAuthorization{
		URL:           authorizationURL,
		BrowserOpened: opened,
		Complete: func(
			ctx context.Context,
			code string,
		) (ports.LinkedIdentity, error) {
			return l.exchange(
				ctx,
				strings.TrimSpace(code),
				verifier,
			)
		},
	}, nil
}

func (l *Linker) exchange(
	ctx context.Context,
	code string,
	verifier string,
) (ports.LinkedIdentity, error) {
	if code == "" {
		return ports.LinkedIdentity{}, errors.New("enter the Dropbox authorization code")
	}

	parameters := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {l.appKey},
		"code_verifier": {verifier},
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		l.tokenURL,
		strings.NewReader(parameters.Encode()),
	)
	if err != nil {
		return ports.LinkedIdentity{}, errors.New("could not prepare Dropbox sign-in")
	}
	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	response, err := l.client.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return ports.LinkedIdentity{}, ctx.Err()
		}
		return ports.LinkedIdentity{}, errors.New("could not connect to Dropbox; try again")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ports.LinkedIdentity{}, fmt.Errorf(
			"Dropbox rejected sign-in (HTTP %d); try a new authorization code",
			response.StatusCode,
		)
	}

	var result struct {
		AccountID    string      `json:"account_id"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		ExpiresIn    json.Number `json:"expires_in"`
	}
	if err := json.NewDecoder(io.LimitReader(
		response.Body,
		1<<20,
	)).Decode(&result); err != nil {
		return ports.LinkedIdentity{}, errors.New("Dropbox returned an invalid sign-in response")
	}
	seconds, err := result.ExpiresIn.Int64()
	if err != nil || seconds <= 0 || seconds > 365*24*60*60 || result.AccessToken == "" || result.RefreshToken == "" {
		return ports.LinkedIdentity{}, errors.New("Dropbox returned incomplete sign-in credentials")
	}
	account, err := identity.NewIdentity(
		result.AccountID,
		identity.ProviderDropbox,
	)
	if err != nil {
		return ports.LinkedIdentity{}, errors.New("Dropbox did not identify the linked account")
	}
	if err := ctx.Err(); err != nil {
		return ports.LinkedIdentity{}, err
	}

	updated := time.Now().UTC()
	credentials, err := identity.NewIdentityCredentials(
		account.ID,
		result.AccessToken,
		result.RefreshToken,
		updated.Add(time.Duration(seconds)*time.Second),
		updated,
	)
	if err != nil {
		return ports.LinkedIdentity{}, err
	}

	return ports.LinkedIdentity{Identity: account, Credentials: credentials}, nil
}

var _ ports.IdentityLinker = (*Linker)(nil)

// END_BLOCK_DROPBOX_AUTHORIZATION
