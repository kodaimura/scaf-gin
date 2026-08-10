package core

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"scaf-gin/config"
)

// NewGormDB initializes a GORM database connection based on configuration.
func NewGormDB(cfg config.Config, log Logger) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	switch cfg.DBEngine {
	case "postgres":
		db, err = gorm.Open(postgres.Open(buildPostgresDSN(cfg)), gormConfig)
	case "mysql":
		db, err = gorm.Open(mysql.Open(buildMySQLDSN(cfg)), gormConfig)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(buildSQLiteDSN(cfg)), gormConfig)
	default:
		return nil, fmt.Errorf("invalid DB_ENGINE. Please choose 'postgres', 'mysql', or 'sqlite3'")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic DB object: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info("database connected via gorm")
	return db, nil
}

func buildPostgresDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)
}

func buildMySQLDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}

func buildSQLiteDSN(cfg config.Config) string {
	return fmt.Sprintf("%s.db", cfg.DBName)
}

func HandleGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrConflict
	}
	if strings.Contains(err.Error(), "SQLSTATE 23505") {
		return ErrConflict
	}

	return NewUnexpectedError(err)
}
