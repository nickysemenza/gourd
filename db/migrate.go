package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func AutoMigrate(_ context.Context, dbConn *sql.DB, up bool) error {
	//nolint: contextcheck
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}
	viper.SetDefault("MIGRATIONS_DIR", "file://./db/migrations")
	viper.AutomaticEnv()
	dir := viper.GetString("MIGRATIONS_DIR")
	log.Infof("Migrations dir: %s", dir)
	m, err := migrate.NewWithDatabaseInstance(
		dir,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	action := m.Down
	if up {
		action = m.Up
	}
	if err := action(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}
	log.Info("db: migrated")
	return nil
}
