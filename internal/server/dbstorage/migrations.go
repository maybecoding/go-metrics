package dbstorage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	//_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//import "github.com/golang-migrate/migrate/v4"

////go:embed migrations/*.sql
//var migrationsDir embed.FS

func runMigrations(dsn string) error {
	//d, err := iofs.New(migrationsDir, "migrations")
	//if err != nil {
	//	return fmt.Errorf("failed to return an iofs driver: %w", err)
	//}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to db for run migrations: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations to DB")
	}
	return nil
}
