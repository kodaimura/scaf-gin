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
	BasicAuthUser string
	BasicAuthPass string
)

var (
	JwtExpiresSeconds string
	JwtSecretKey string
)

var (
	LogLevel  string
)

func init() {
	env := os.Getenv("ENV")
	if env != "" {
		env = "." + env
	}
	fmt.Println(fmt.Sprintf("config/env/.env%s", env))
	err := godotenv.Load(fmt.Sprintf("config/env/.env%s", env))
	if err != nil {
		log.Panic("Failed to load environment variables:", err)
	}

	AppName = getEnv("APP_NAME")
	AppHost = getEnv("APP_HOST", "localhost")
	AppPort = getEnv("APP_PORT", "3000")

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

	BasicAuthUser = getEnv("BASIC_AUTH_USER")
	BasicAuthPass = getEnv("BASIC_AUTH_PASSWORD")

	JwtExpiresSeconds = getEnv("JWT_EXPIRES_SECONDS", "3600")
	JwtSecretKey = getEnv("JWT_SECRET_KEY", "secret")

	LogLevel  = getEnv("LOG_LEVEL", "INFO")
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}