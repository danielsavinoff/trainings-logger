// START_MODULE_CONTRACT
//   PURPOSE: Transfer workout backups through a linked Dropbox account.
//   SCOPE: Local credential access, token refresh, bounded upload/download, sanitized errors.
//   DEPENDS: M-DOMAIN-IDENTITY, M-DOMAIN-STATE, M-APP-REMOTE-STATE, M-INFRA-DROPBOX-AUTH
//   LINKS: M-INFRA-DROPBOX-STATE, V-M-INFRA-DROPBOX-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   StateStore - retain configured Dropbox and credential ports
//   NewStateStore - wire the real content API
//   Write - validate and upload a versioned backup
//   Read - download, decode, and validate a bounded backup
//   contentRequest - perform content requests with one expired-token retry
//   contentError - classify provider errors without exposing response bodies
//   dropboxRequestID - retain only bounded hexadecimal provider request IDs
//   accessToken - select the latest linked account and refresh when needed
//   refresh - exchange and persist replacement credentials
// END_MODULE_MAP

package dropbox

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
)

// START_BLOCK_DROPBOX_SNAPSHOTS
const backupPath = "/trainings-logger/data.json"
const maxSnapshotBytes = 16 << 20

type StateStore struct {
	linker      *Linker
	accounts    identity.IdentityRepository
	credentials identity.IdentityCredentialsRepository
	contentURL  string
}

func NewStateStore(
	linker *Linker,
	accounts identity.IdentityRepository,
	credentials identity.IdentityCredentialsRepository,
) *StateStore {
	return &StateStore{linker: linker, accounts: accounts, credentials: credentials, contentURL: "https://content.dropboxapi.com/2/files/"}
}

func (s *StateStore) Write(
	ctx context.Context,
	snapshot state.Snapshot,
) error {
	validated, err := state.NewSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sets,
		snapshot.Preferences,
	)
	if err != nil {
		return err
	}
	data, err := json.Marshal(validated)
	if err != nil {
		return errors.New("could not encode the workout backup")
	}
	if len(data) > maxSnapshotBytes {
		return errors.New("workout backup exceeds the 16 MiB transfer limit")
	}
	argument := `{"path":"` + backupPath + `","mode":"overwrite","autorename":false,"mute":true}`
	response, err := s.contentRequest(
		ctx,
		"upload",
		argument,
		data,
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(
		io.Discard,
		io.LimitReader(
			response.Body,
			1<<20,
		),
	)
	if err != nil {
		return errors.New("Dropbox upload response could not be read; try saving again")
	}
	return nil
}

func (s *StateStore) Read(ctx context.Context) (state.Snapshot, error) {
	response, err := s.contentRequest(
		ctx,
		"download",
		`{"path":"`+backupPath+`"}`,
		nil,
	)
	if err != nil {
		return state.Snapshot{}, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(
		response.Body,
		maxSnapshotBytes+1,
	))
	if err != nil {
		return state.Snapshot{}, errors.New("Dropbox backup could not be downloaded")
	}
	if len(data) > maxSnapshotBytes {
		return state.Snapshot{}, errors.New("Dropbox backup exceeds the 16 MiB transfer limit")
	}
	var snapshot state.Snapshot
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snapshot); err != nil {
		return state.Snapshot{}, errors.New("Dropbox backup is not a valid workout snapshot")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return state.Snapshot{}, errors.New("Dropbox backup contains trailing data")
	}
	return state.NewSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sets,
		snapshot.Preferences,
	)
}

func (s *StateStore) contentRequest(
	ctx context.Context,
	operation string,
	argument string,
	data []byte,
) (*http.Response, error) {
	for attempt := 0; attempt < 2; attempt++ {
		token, err := s.accessToken(
			ctx,
			attempt == 1,
		)
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.contentURL+operation,
			bytes.NewReader(data),
		)
		if err != nil {
			return nil, errors.New("could not prepare the Dropbox transfer")
		}
		request.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
		request.Header.Set(
			"Dropbox-API-Arg",
			argument,
		)
		if operation == "upload" {
			request.Header.Set(
				"Content-Type",
				"application/octet-stream",
			)
		}
		response, err := s.linker.client.Do(request)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, errors.New("could not connect to Dropbox; try again")
		}
		if response.StatusCode == http.StatusOK {
			return response, nil
		}
		err, expired := contentError(
			response,
			operation,
		)
		response.Body.Close()
		if !expired || attempt == 1 {
			return nil, err
		}
	}
	return nil, errors.New("Dropbox authorization expired; press r to link again")
}

func contentError(
	response *http.Response,
	operation string,
) (error, bool) {
	data, _ := io.ReadAll(io.LimitReader(
		response.Body,
		1<<20,
	))
	var result struct {
		Summary string `json:"error_summary"`
		Error   struct {
			Tag           string `json:".tag"`
			RequiredScope string `json:"required_scope"`
		} `json:"error"`
	}
	_ = json.Unmarshal(
		data,
		&result,
	)
	if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusUnauthorized {
		missingScope := strings.HasPrefix(
			result.Summary,
			"missing_scope/",
		) || result.Error.Tag == "missing_scope" || strings.Contains(
			string(data),
			"does not have the required scope",
		) || strings.Contains(
			string(data),
			"Your app is not authorized to use this endpoint",
		)
		if missingScope {
			requiredScope := ""
			for _, scope := range []string{"files.content.read", "files.content.write"} {
				if result.Error.RequiredScope == scope || strings.Contains(
					string(data),
					scope,
				) {
					requiredScope = scope
					break
				}
			}
			if requiredScope == "" && strings.Contains(
				string(data),
				"Your app is not authorized to use this endpoint",
			) {
				switch operation {
				case "upload":
					requiredScope = "files.content.write"
				case "download":
					requiredScope = "files.content.read"
				}
			}
			if requiredScope != "" {
				return fmt.Errorf(
					"Dropbox requires %s; enable it in App Console > Permissions, Submit, then press r to link again",
					requiredScope,
				), false
			}
			return errors.New("Dropbox file permissions are missing; enable file read/write and press r to link again"), false
		}
	}
	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("Dropbox authorization is no longer valid; press r to link again"), strings.HasPrefix(
			result.Summary,
			"expired_access_token/",
		)
	}
	if operation == "download" && response.StatusCode == http.StatusConflict && strings.HasPrefix(
		result.Summary,
		"path/not_found/",
	) {
		return errors.New("no Dropbox backup found; press s to save one first"), false
	}
	if response.StatusCode == http.StatusTooManyRequests {
		return errors.New("Dropbox is busy; try again later"), false
	}
	if response.StatusCode == http.StatusBadRequest {
		requestID := dropboxRequestID(response.Header)
		if requestID != "" {
			return fmt.Errorf(
				"Dropbox %s rejected the request (HTTP 400; request %s)",
				operation,
				requestID,
			), false
		}
		return fmt.Errorf(
			"Dropbox %s rejected the request (HTTP 400)",
			operation,
		), false
	}
	return fmt.Errorf(
		"Dropbox %s failed (HTTP %d); try again",
		operation,
		response.StatusCode,
	), false
}

func dropboxRequestID(header http.Header) string {
	value := header.Get("X-Dropbox-Request-Id")
	if len(value) == 0 || len(value) > 64 {
		return ""
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') && (character < 'A' || character > 'F') {
			return ""
		}
	}
	return value
}

// END_BLOCK_DROPBOX_SNAPSHOTS

// START_BLOCK_DROPBOX_TOKEN_REFRESH
func (s *StateStore) accessToken(
	ctx context.Context,
	forceRefresh bool,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	credentials, err := s.credentials.FindLatest(ctx)
	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return "", errors.New("no Dropbox account linked; press r to set up Dropbox")
	}
	if err != nil {
		return "", errors.New("could not read the linked Dropbox account")
	}
	account, err := s.accounts.FindByID(
		ctx,
		credentials.IdentityID,
	)
	if err != nil || account.Provider != identity.ProviderDropbox {
		return "", errors.New("the linked Dropbox account is unavailable; press r to link again")
	}
	if !forceRefresh && credentials.AccessTokenExpiresAt.After(time.Now().Add(time.Minute)) {
		return credentials.AccessToken, nil
	}

	return s.refresh(
		ctx,
		credentials,
	)
}

func (s *StateStore) refresh(
	ctx context.Context,
	previous identity.IdentityCredentials,
) (string, error) {
	if s.linker.appKey == "" {
		return "", errors.New("Dropbox app key is not configured")
	}
	parameters := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {previous.RefreshToken},
		"client_id":     {s.linker.appKey},
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.linker.tokenURL,
		strings.NewReader(parameters.Encode()),
	)
	if err != nil {
		return "", errors.New("could not prepare Dropbox authorization refresh")
	}
	request.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	response, err := s.linker.client.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", errors.New("could not refresh Dropbox authorization; try again")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Dropbox authorization could not be refreshed; press r to link again")
	}
	var value struct {
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		AccountID    string      `json:"account_id"`
		ExpiresIn    json.Number `json:"expires_in"`
	}
	if err := json.NewDecoder(io.LimitReader(
		response.Body,
		1<<20,
	)).Decode(&value); err != nil {
		return "", errors.New("Dropbox returned invalid refreshed credentials")
	}
	seconds, err := value.ExpiresIn.Int64()
	if err != nil || seconds <= 0 || seconds > 365*24*60*60 || value.AccessToken == "" || (value.AccountID != "" && value.AccountID != previous.IdentityID) {
		return "", errors.New("Dropbox returned incomplete refreshed credentials")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	refresh := previous.RefreshToken
	if value.RefreshToken != "" {
		refresh = value.RefreshToken
	}

	now := time.Now().UTC()
	credentials, err := identity.NewIdentityCredentials(
		previous.IdentityID,
		value.AccessToken,
		refresh,
		now.Add(time.Duration(seconds)*time.Second),
		now,
	)
	if err != nil {
		return "", err
	}
	if err := s.credentials.Save(
		ctx,
		credentials,
	); err != nil {
		return "", errors.New("could not save refreshed Dropbox authorization")
	}
	return value.AccessToken, nil
}

var _ ports.BackupWriter = (*StateStore)(nil)
var _ ports.BackupReader = (*StateStore)(nil)

// END_BLOCK_DROPBOX_TOKEN_REFRESH
