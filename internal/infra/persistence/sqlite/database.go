package sqlite

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(ctx context.Context) (*sql.DB, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return OpenFile(ctx, filepath.Join(dir, "trainlog", "trainlog.db"))
}

func OpenFile(ctx context.Context, path string) (*sql.DB, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	dsn := (&url.URL{Scheme: "file", Path: filepath.ToSlash(path)}).String() +
		"?_pragma=foreign_keys(1)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Serialize SQLite writes and keep connection PRAGMAs stable.
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, err
	}

	if err := Migrate(db); err != nil {
		db.Close()

		return nil, err
	}

	return db, nil
}
