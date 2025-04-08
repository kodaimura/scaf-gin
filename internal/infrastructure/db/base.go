package db

import (
	"fmt"
	
	"goscaf/config"
)

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