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

func NewGormDB() *gorm.DB {
	dbEngine := config.DBEngine
	dbName := config.DBName
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPass := config.DBPass

	var db *gorm.DB
	var err error

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	switch dbEngine {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
		fmt.Println(dsn)
		db, err = gorm.Open(postgres.Open(dsn), config)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		db, err = gorm.Open(mysql.Open(dsn), config)
	case "sqlite3":
		dsn := fmt.Sprintf("%s.db", dbName)
		db, err = gorm.Open(sqlite.Open(dsn), config)
	default:
		log.Panic("Error: must specify a valid DB_DRIVER: 'postgres', 'mysql', or 'sqlite3'.")
	}

	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Panic(err)
	}

	log.Println("Successfully connected to GORM database")

	return db
}