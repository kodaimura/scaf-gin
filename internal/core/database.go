package core

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"scaf-gin/config"
)

// NewGormDB initializes a GORM database connection based on configuration.
func NewGormDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	switch config.DBEngine {
	case "postgres":
		db, err = gorm.Open(postgres.Open(buildPostgresDSN()), gormConfig)
	case "mysql":
		db, err = gorm.Open(mysql.Open(buildMySQLDSN()), gormConfig)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(buildSQLiteDSN()), gormConfig)
	default:
		failDB("invalid DB_ENGINE. Please choose 'postgres', 'mysql', or 'sqlite3'.")
	}

	if err != nil {
		failDB("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		failDB("failed to get generic DB object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		failDB("database ping failed: %v", err)
	}

	Logger.Info("database connected via gorm")
	return db
}

// NewSqlxDB initializes a SQLx database connection based on configuration.
func NewSqlxDB() *sqlx.DB {
	var (
		db  *sqlx.DB
		err error
	)

	switch config.DBEngine {
	case "postgres":
		db, err = sqlx.Connect("postgres", buildPostgresDSN())
	case "mysql":
		db, err = sqlx.Connect("mysql", buildMySQLDSN())
	case "sqlite3":
		db, err = sqlx.Connect("sqlite3", buildSQLiteDSN())
	default:
		failDB("invalid DB_ENGINE. Please choose 'postgres', 'mysql', or 'sqlite3'.")
	}

	if err != nil {
		failDB("failed to connect using sqlx: %v", err)
	}

	if err := db.Ping(); err != nil {
		failDB("database ping failed: %v", err)
	}

	Logger.Info("database connected via sqlx")
	return db
}

func buildPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPass, config.DBName, config.DBPort,
	)
}

func buildMySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName,
	)
}

func buildSQLiteDSN() string {
	return fmt.Sprintf("%s.db", config.DBName)
}

func failDB(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Logger.Error("%s", message)
	panic(message)
}
