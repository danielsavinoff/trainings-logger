// START_MODULE_CONTRACT
//   PURPOSE: Protect Dropbox tokens before they reach persistence.
//   SCOPE: Private local encryption keys and authenticated token encryption.
//   DEPENDS: none
//   LINKS: M-INFRA-DROPBOX-AUTH, V-M-INFRA-DROPBOX-AUTH
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   credentialCipher - load or create a private AES-GCM key
//   encryptToken - encrypt one token with account and role binding
//   decryptToken - authenticate and decode a stored token
// END_MODULE_MAP

package dropbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// START_BLOCK_CREDENTIAL_ENCRYPTION
func credentialCipher(path string) (cipher.AEAD, error) {
	if path == "" {
		return nil, errors.New("Dropbox credential storage is not configured")
	}

	key, err := os.ReadFile(path)
	if errors.Is(
		err,
		os.ErrNotExist,
	) {
		if err := os.MkdirAll(
			filepath.Dir(path),
			0700,
		); err != nil {
			return nil, errors.New("could not prepare Dropbox credential storage")
		}

		key = make(
			[]byte,
			32,
		)
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		file, err := os.OpenFile(
			path,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL,
			0600,
		)
		if errors.Is(
			err,
			os.ErrExist,
		) {
			return credentialCipher(path)
		}
		if err != nil {
			return nil, errors.New("could not create Dropbox credential key")
		}

		_, writeErr := file.Write(key)
		closeErr := file.Close()
		if writeErr != nil || closeErr != nil {
			return nil, errors.New("could not save Dropbox credential key")
		}
	} else if err != nil {
		return nil, errors.New("could not read Dropbox credential key")
	}

	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 || len(key) != 32 {
		return nil, errors.New("Dropbox credential key must be a private 32-byte file")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}

func encryptToken(
	key cipher.AEAD,
	token string,
	binding string,
) (string, error) {
	nonce := make(
		[]byte,
		key.NonceSize(),
	)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	encrypted := key.Seal(
		nonce,
		nonce,
		[]byte(token),
		[]byte(binding),
	)

	return "v1:" + base64.RawStdEncoding.EncodeToString(encrypted), nil
}

func decryptToken(
	key cipher.AEAD,
	value string,
	binding string,
) (string, error) {
	if !strings.HasPrefix(
		value,
		"v1:",
	) {
		return "", errors.New("unsupported Dropbox credential encoding")
	}

	encrypted, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(
		value,
		"v1:",
	))
	if err != nil || len(encrypted) < key.NonceSize() {
		return "", errors.New("invalid encrypted Dropbox credential")
	}

	plain, err := key.Open(
		nil,
		encrypted[:key.NonceSize()],
		encrypted[key.NonceSize():],
		[]byte(binding),
	)
	if err != nil {
		return "", errors.New("could not authenticate Dropbox credential")
	}

	return string(plain), nil
}

// END_BLOCK_CREDENTIAL_ENCRYPTION
