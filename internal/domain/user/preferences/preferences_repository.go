package preferences

import "context"

type PreferencesRepository interface {
	Save(ctx context.Context, preferences Preferences) error
	Find(ctx context.Context) (Preferences, error)
}
