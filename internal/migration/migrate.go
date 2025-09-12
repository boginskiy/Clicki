package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const MigrationsDir = "./migrations"

func ApplyMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatalf("Failed to create driver instance: %v\\n", err)
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filepath.Join(os.Getenv("PWD"), MigrationsDir)),
		"postgres", driver)

	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v\\n", err)
		return err
	}
	defer func(mgr *migrate.Migrate) { mgr.Close() }(migrator)

	// ctx := context.Background()

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v\\n", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
