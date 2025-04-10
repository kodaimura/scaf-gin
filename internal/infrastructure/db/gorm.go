package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"goscaf/config"
)

// NewGormDB initializes a GORM database connection based on configuration.
// Supports postgres, mysql, and sqlite3.
func NewGormDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Prevent pluralized table names
		},
	}

	switch config.DBEngine {
	case "postgres":
		dsn := internal/infrastructure/db/base.go
		db, err = gorm.Open(postgres.Open(buildPostgresDSN()), gormConfig)
	case "mysql":
		db, err = gorm.Open(mysql.Open(buildMySQLDSN()), gormConfig)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(buildSQLiteDSN()), gormConfig)
	default:
		log.Panic("❌ Invalid DB_ENGINE. Please choose 'postgres', 'mysql', or 'sqlite3'.")
	}

	if err != nil {
		log.Panicf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("❌ Failed to get generic DB object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Panicf("❌ Database ping failed: %v", err)
	}

	log.Println("✅ Successfully connected to the database via GORM.")
	return db
}


// ===============================
// Common for the db package.
// ===============================
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