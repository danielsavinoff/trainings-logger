// START_MODULE_CONTRACT
//   PURPOSE: Represent a numbered workout session.
//   SCOPE: Validate UUID identity, index, and UTC creation time.
//   DEPENDS: none
//   LINKS: M-DOMAIN-SESSION, V-M-DOMAIN-SESSION
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Session - retain validated session identity, index and creation time
//   NewSession - validate and construct a numbered session
// END_MODULE_MAP

package session

import (
	"time"
	"uuid"
)

// START_BLOCK_SESSION_MODEL
type Session struct {
	ID        string    `json:"id"`
	Index     int       `json:"index"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSession(
	id string,
	index int,
	createdAt *time.Time,
) (Session, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Session{}, ErrInvalidSessionID
	}

	if index < 1 {
		return Session{}, ErrInvalidIndex
	}

	now := time.Now().UTC()
	if createdAt != nil {
		now = createdAt.UTC()
	}

	if now.IsZero() || now.Year() < 1 || now.Year() > 9999 {
		return Session{}, ErrInvalidCreatedAt
	}

	return Session{
		ID:        id,
		Index:     index,
		CreatedAt: now,
	}, nil
}

// END_BLOCK_SESSION_MODEL
