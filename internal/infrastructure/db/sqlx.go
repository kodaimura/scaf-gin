package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"scaf-gin/config"
)

// NewSqlxDB initializes a SQLx database connection based on configuration.
// Supports 'postgres', 'mysql', and 'sqlite3'.
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
		log.Panic("❌ Invalid DB_ENGINE. Choose 'postgres', 'mysql', or 'sqlite3'.")
	}

	if err != nil {
		log.Panicf("❌ Failed to connect using sqlx: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Panicf("❌ Database ping failed: %v", err)
	}

	log.Printf("✅ Successfully connected to %s via sqlx\n", config.DBEngine)
	return db
}