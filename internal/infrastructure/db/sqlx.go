package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"goscaf/config"
)

func NewSqlxDB() (*sqlx.DB, error) {
	dbEngine := config.DBEngine
	dbName := config.DBName
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPass := config.DBPass

	var db *sqlx.DB
	var err error

	switch dbEngine {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
		db, err = sqlx.Connect("postgres", dsn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		db, err = sqlx.Connect("mysql", dsn)
	case "sqlite3":
		dsn := fmt.Sprintf("%s.db", dbName)
		db, err = sqlx.Connect("sqlite3", dsn)
	default:
		return nil, fmt.Errorf("unsupported DB_ENGINE: %s", dbEngine)
	}

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to SQLx database")

	return db, nil
}
