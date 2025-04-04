package config

import (
	"os"
	"log"
	"fmt"

	"github.com/joho/godotenv"
)

var (
	AppName string
	AppHost string
	AppPort string
)

var (
	DBEngine string
	DBName   string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
)

var (
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
)

var (
	AuthUser string
	AuthPass string
)

var (
	JWTSecret string
	LogLevel  string
)

func init() {
	env := os.Getenv("ENV")
	if env != "" {
		env = "." + env
	}
	err := godotenv.Load(fmt.Sprintf("config/env/.env%s", env))
	if err != nil {
		log.Panic("Failed to load environment variables:", err)
	}

	AppName = os.Getenv("APP_NAME")
	AppHost = os.Getenv("APP_HOST")
	AppPort = os.Getenv("APP_PORT")

	DBEngine = os.Getenv("DB_ENGINE")
	DBName = os.Getenv("DB_NAME")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASSWORD")

	SMTPHost = os.Getenv("SMTP_HOST")
	SMTPPort = os.Getenv("SMTP_PORT")
	SMTPUser = os.Getenv("SMTP_USER")
	SMTPPass = os.Getenv("SMTP_PASSWORD")

	AuthUser = os.Getenv("BASIC_AUTH_USER")
	AuthPass = os.Getenv("BASIC_AUTH_PASSWORD")

	JWTSecret = os.Getenv("JWT_SECRET")
	LogLevel  = os.Getenv("LOG_LEVEL")
}