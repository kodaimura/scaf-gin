package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
	MailFrom string
)

var (
	BasicAuthUser string
	BasicAuthPass string
)

var (
	JwtSecretKey       string
	JwtRefreshSecretKey string
	AuthExpiresSeconds int
	RefreshExpiresSeconds int
)

var (
	SecureCookie   bool
	LogLevel       string
	FrontendOrigin string
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

	AppName = getEnv("APP_NAME")
	AppHost = getEnv("APP_HOST", "localhost")
	AppPort = getEnv("APP_PORT", "8000")

	DBEngine = getEnv("DB_ENGINE")
	DBName = getEnv("DB_NAME")
	DBHost = getEnv("DB_HOST")
	DBPort = getEnv("DB_PORT")
	DBUser = getEnv("DB_USER")
	DBPass = getEnv("DB_PASSWORD")

	SMTPHost = getEnv("SMTP_HOST")
	SMTPPort = getEnv("SMTP_PORT")
	SMTPUser = getEnv("SMTP_USER")
	SMTPPass = getEnv("SMTP_PASSWORD")
	MailFrom = getEnv("MAIL_FROM")

	BasicAuthUser = getEnv("BASIC_AUTH_USER")
	BasicAuthPass = getEnv("BASIC_AUTH_PASSWORD")

	JwtSecretKey = getEnv("JWT_SECRET_KEY", "secret")
	AuthExpiresSeconds, err = strconv.Atoi(getEnv("AUTH_EXPIRES_SECONDS", "3600"))
	if err != nil {
		log.Fatalf("unable to convert AUTH_EXPIRES_SECONDS from environment to integer: %v", err)
	}

	JwtRefreshSecretKey = getEnv("JWT_REFRESH_SECRET_KEY", "secret")
	RefreshExpiresSeconds, err = strconv.Atoi(getEnv("REFRESH_EXPIRES_SECONDS", "86400"))
	if err != nil {
		log.Fatalf("unable to convert REFRESH_EXPIRES_SECONDS from environment to integer: %v", err)
	}

	SecureCookie, err = strconv.ParseBool(getEnv("SECURE_COOKIE", "false"))
	if err != nil {
		log.Fatalf("unable to convert SECURE_COOKIE from environment to boolean: %v", err)
	}
	LogLevel = getEnv("LOG_LEVEL", "INFO")
	FrontendOrigin = getEnv("FRONTEND_ORIGIN", "http://localhost:3000")
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
