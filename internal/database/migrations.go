package database

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunDBMigrations(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		slog.Error("cannot create new migrate instance", "err", err)
		return
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrate up", "err", err)
		return
	}

	slog.Info("DB migrated successfully")
}
