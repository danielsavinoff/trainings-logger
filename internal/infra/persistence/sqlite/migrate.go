package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/danielsavinoff/trainings-logger/migrations"
	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(db *sql.DB) error {
	source, err := iofs.New(migrations.Files, "sqlite")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}
	defer source.Close()

	driver, err := migratesqlite.WithInstance(db, &migratesqlite.Config{})
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	runner, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	// The caller owns db; runner.Close would also close that shared connection pool.
	if err := runner.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply SQLite migrations: %w", err)
	}

	return nil
}
