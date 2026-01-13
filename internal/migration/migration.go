package migration

import (
	"database/sql"
	"fmt"
	"log"
)

type Migration struct {
	Version uint64
	Name string
	Up string
	Down string
}

var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_users_table",
		Up: `
			CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,
				email VARCHAR(255) NOT NULL,
				username VARCHAR(255) NOT NULL,
				password TEXT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				deleted_at TIMESTAMPTZ
			);

			CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email_active
				ON users(email) WHERE deleted_at IS NULL;

			CREATE UNIQUE INDEX IF NOT EXISTS uq_users_username_active
				ON users(username) WHERE deleted_at IS NULL;

			CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

			CREATE OR REPLACE FUNCTION set_updated_at()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW.updated_at = now();
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;

			DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
			CREATE TRIGGER trg_users_set_updated_at
				BEFORE UPDATE ON users
				FOR EACH ROW
				EXECUTE FUNCTION set_updated_at();
		`,
		Down: `
			DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;
			DROP FUNCTION IF EXISTS set_updated_at;

			DROP INDEX IF EXISTS idx_users_deleted_at;
			DROP INDEX IF EXISTS uq_users_email_active;
			DROP INDEX IF EXISTS uq_users_username_active;

			DROP TABLE IF EXISTS users CASCADE;
		`,
	},
	{
		Version: 2,
		Name:    "create_notes_table",
		Up: `
			CREATE TABLE IF NOT EXISTS notes (
				id BIGSERIAL PRIMARY KEY,
				user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				title VARCHAR(255) NOT NULL,
				content TEXT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				deleted_at TIMESTAMPTZ
			);

			CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
			CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
			CREATE INDEX IF NOT EXISTS idx_notes_deleted_at ON notes(deleted_at);

			CREATE INDEX IF NOT EXISTS idx_notes_user_created_active
				ON notes(user_id, created_at DESC)
				WHERE deleted_at IS NULL;

			DROP TRIGGER IF EXISTS trg_notes_set_updated_at ON notes;
			CREATE TRIGGER trg_notes_set_updated_at
				BEFORE UPDATE ON notes
				FOR EACH ROW
				EXECUTE FUNCTION set_updated_at();
		`,
		Down: `
			DROP TRIGGER IF EXISTS trg_notes_set_updated_at ON notes;

			DROP INDEX IF EXISTS idx_notes_user_created_active;
			DROP INDEX IF EXISTS idx_notes_user_id;
			DROP INDEX IF EXISTS idx_notes_created_at;
			DROP INDEX IF EXISTS idx_notes_deleted_at;

			DROP TABLE IF EXISTS notes CASCADE;
		`,
	},
	{
		Version: 3,
		Name:    "create_tags_table",
		Up: `
			CREATE TABLE IF NOT EXISTS tags (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL UNIQUE,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);

			CREATE TABLE IF NOT EXISTS note_tags (
				note_id BIGINT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
				tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
				PRIMARY KEY (note_id, tag_id)
			);

			CREATE INDEX IF NOT EXISTS idx_note_tags_note_id ON note_tags(note_id);
			CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags(tag_id);
			CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
		`,
		Down: `
			DROP INDEX IF EXISTS idx_tags_name;
			DROP INDEX IF EXISTS idx_note_tags_note_id;
			DROP INDEX IF EXISTS idx_note_tags_tag_id;

			DROP TABLE IF EXISTS note_tags;
			DROP TABLE IF EXISTS tags CASCADE;
		`,
	},
}


func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func MigrateUp(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	log.Printf("Current migration version: %d", currentVersion)

	for _, migration := range migrations {
		if migration.Version <= uint64(currentVersion) {
			continue
		}
		log.Printf("Running migration %d: %s", migration.Version, migration.Name)

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Run migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}
		log.Printf("Migration %d completed successfully", migration.Version)
	}
	log.Println("All migrations completed!")
	return nil
}

func MigrateDown(db *sql.DB) error {
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if currentVersion == 0 {
		log.Println("No migrations to roll back")
		return nil
	}
		// Find the migration to roll back
	var targetMigration *Migration
	for i := range migrations {
		if migrations[i].Version == uint64(currentVersion) {
			targetMigration = &migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration version %d not found", currentVersion)
	}

	log.Printf("Rolling back migration %d: %s", targetMigration.Version, targetMigration.Name)
		// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// Run down migration
	if _, err := tx.Exec(targetMigration.Down); err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback %d failed: %w", targetMigration.Version, err)
	}
	// Remove migration record
	if _, err := tx.Exec(
		"DELETE FROM schema_migrations WHERE version = $1",
		targetMigration.Version,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", targetMigration.Version, err)
	}
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %w", targetMigration.Version, err)
	}
	log.Printf("Migration %d rolled back successfully", targetMigration.Version)
	return nil
}
