package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"goscaf/config"
)

func NewGormDB() (*gorm.DB, error) {
	dbEngine := config.DBEngine
	dbName := config.DBName
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPass := config.DBPass

	var db *gorm.DB
	var err error

	// DB_ENGINEによって接続先を選択
	switch dbEngine {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite3":
		dsn := fmt.Sprintf("%s.db", dbName)
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported DB_ENGINE: %s", dbEngine)
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to GORM database")

	return db, nil
}