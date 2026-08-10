package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"scaf-gin/config"
	"scaf-gin/internal/core"

	"gorm.io/gorm"
)

type migration struct {
	Version int
	Name    string
	Path    string
}

func main() {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cfg := config.Current
	log := core.NewJSONLogger(cfg.LogLevel, "migrate")
	db, err := core.NewGormDB(cfg, log)
	if err != nil {
		log.Error("database connection failed: %v", err)
		os.Exit(1)
	}

	switch command {
	case "up":
		err = migrateUp(db, log)
	case "current":
		err = printCurrent(db)
	case "history":
		err = printHistory(db)
	default:
		err = fmt.Errorf("unknown migrate command: %s", command)
	}
	if err != nil {
		log.Error("migration failed: %v", err)
		os.Exit(1)
	}
}

func migrateUp(db *gorm.DB, log core.Logger) error {
	migrations, err := loadMigrations("migrations")
	if err != nil {
		return err
	}

	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	applied, err := appliedVersions(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		sql, err := os.ReadFile(migration.Path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", migration.Name, err)
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(sql)).Error; err != nil {
				return fmt.Errorf("apply migration %s: %w", migration.Name, err)
			}
			return tx.Exec(
				"INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)",
				migration.Version,
				migration.Name,
				time.Now().UTC(),
			).Error
		})
		if err != nil {
			return err
		}

		log.Info("applied migration: version=%03d name=%s", migration.Version, migration.Name)
	}

	log.Info("database migrations are up to date")
	return nil
}

func printCurrent(db *gorm.DB) error {
	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	var current appliedMigration
	err := db.Raw(
		"SELECT version, name, applied_at FROM schema_migrations ORDER BY version DESC LIMIT 1",
	).Scan(&current).Error
	if err != nil {
		return err
	}
	if current.Version == 0 {
		fmt.Println("no migrations applied")
		return nil
	}
	fmt.Printf("%03d %s %s\n", current.Version, current.Name, current.AppliedAt.Format(time.RFC3339))
	return nil
}

func printHistory(db *gorm.DB) error {
	if err := ensureMigrationTable(db); err != nil {
		return err
	}

	var history []appliedMigration
	err := db.Raw(
		"SELECT version, name, applied_at FROM schema_migrations ORDER BY version",
	).Scan(&history).Error
	if err != nil {
		return err
	}
	for _, migration := range history {
		fmt.Printf("%03d %s %s\n", migration.Version, migration.Name, migration.AppliedAt.Format(time.RFC3339))
	}
	return nil
}

func ensureMigrationTable(db *gorm.DB) error {
	return db.Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)`).Error
}

func appliedVersions(db *gorm.DB) (map[int]bool, error) {
	var applied []appliedMigration
	if err := db.Raw("SELECT version, name, applied_at FROM schema_migrations").Scan(&applied).Error; err != nil {
		return nil, err
	}

	versions := make(map[int]bool, len(applied))
	for _, migration := range applied {
		versions[migration.Version] = true
	}
	return versions, nil
}

func loadMigrations(dir string) ([]migration, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(`^(\d+)_(.+)\.sql$`)
	migrations := make([]migration, 0, len(files))
	for _, path := range files {
		name := filepath.Base(path)
		match := pattern.FindStringSubmatch(name)
		if len(match) != 3 {
			return nil, fmt.Errorf("invalid migration filename: %s", name)
		}

		version, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version: %s", name)
		}

		migrations = append(migrations, migration{
			Version: version,
			Name:    strings.TrimSuffix(name, ".sql"),
			Path:    path,
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for i := 1; i < len(migrations); i++ {
		if migrations[i-1].Version == migrations[i].Version {
			return nil, fmt.Errorf("duplicate migration version: %03d", migrations[i].Version)
		}
	}

	return migrations, nil
}

type appliedMigration struct {
	Version   int
	Name      string
	AppliedAt time.Time
}
