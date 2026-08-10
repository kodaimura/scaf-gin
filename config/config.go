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
	EnableSignup    bool
	AuthLoginIDMode string
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
	MailProvider string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	MailFrom     string
)

var (
	BasicAuthUser string
	BasicAuthPass string
)

var (
	AccessTokenSecret                    string
	RefreshTokenSecret                   string
	AccessTokenExpiresSeconds            int
	RefreshTokenExpiresSeconds           int
	RefreshTokenRememberMeExpiresSeconds int

	CookieAccessSecure    bool
	CookieRefreshSecure   bool
	CookieAccessHttpOnly  bool
	CookieRefreshHttpOnly bool
)

var (
	PasswordResetURLBase               string
	PasswordResetTokenExpiresMinutes   int
	PasswordResetResendIntervalMinutes int
)

var (
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
	EnableSignup, err = strconv.ParseBool(getEnv("ENABLE_SIGNUP", "true"))
	if err != nil {
		log.Fatalf("unable to convert ENABLE_SIGNUP from environment to boolean: %v", err)
	}
	AuthLoginIDMode = getEnv("AUTH_LOGIN_ID_MODE", "email")

	DBEngine = getEnv("DB_ENGINE")
	DBName = getEnv("DB_NAME")
	DBHost = getEnv("DB_HOST")
	DBPort = getEnv("DB_PORT")
	DBUser = getEnv("DB_USER")
	DBPass = getEnv("DB_PASSWORD")

	MailProvider = getEnv("MAIL_PROVIDER", "mailhog")
	SMTPHost = getEnv("SMTP_HOST")
	SMTPPort = getEnv("SMTP_PORT", "587")
	SMTPUser = getEnv("SMTP_USER")
	SMTPPass = getEnv("SMTP_PASSWORD")
	MailFrom = getEnv("MAIL_FROM", "no-reply@example.local")

	BasicAuthUser = getEnv("BASIC_AUTH_USER")
	BasicAuthPass = getEnv("BASIC_AUTH_PASSWORD")

	AccessTokenSecret = getEnv("ACCESS_TOKEN_SECRET", "secret")
	AccessTokenExpiresSeconds, err = strconv.Atoi(getEnv("ACCESS_TOKEN_EXPIRES_SECONDS", "900"))
	if err != nil {
		log.Fatalf("unable to convert ACCESS_TOKEN_EXPIRES_SECONDS from environment to integer: %v", err)
	}

	RefreshTokenSecret = getEnv("REFRESH_TOKEN_SECRET", "secret")
	RefreshTokenExpiresSeconds, err = strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRES_SECONDS", "2592000"))
	if err != nil {
		log.Fatalf("unable to convert REFRESH_TOKEN_EXPIRES_SECONDS from environment to integer: %v", err)
	}
	RefreshTokenRememberMeExpiresSeconds, err = strconv.Atoi(getEnv("REFRESH_TOKEN_REMEMBER_ME_EXPIRES_SECONDS", "2592000"))
	if err != nil {
		log.Fatalf("unable to convert REFRESH_TOKEN_REMEMBER_ME_EXPIRES_SECONDS from environment to integer: %v", err)
	}

	CookieAccessSecure, err = strconv.ParseBool(getEnv("COOKIE_ACCESS_SECURE", "true"))
	if err != nil {
		log.Fatalf("unable to convert COOKIE_ACCESS_SECURE from environment to boolean: %v", err)
	}

	CookieRefreshSecure, err = strconv.ParseBool(getEnv("COOKIE_REFRESH_SECURE", "true"))
	if err != nil {
		log.Fatalf("unable to convert COOKIE_REFRESH_SECURE from environment to boolean: %v", err)
	}

	CookieAccessHttpOnly, err = strconv.ParseBool(getEnv("COOKIE_ACCESS_HTTPONLY", "true"))
	if err != nil {
		log.Fatalf("unable to convert COOKIE_ACCESS_HTTPONLY from environment to boolean: %v", err)
	}

	CookieRefreshHttpOnly, err = strconv.ParseBool(getEnv("COOKIE_REFRESH_HTTPONLY", "true"))
	if err != nil {
		log.Fatalf("unable to convert COOKIE_REFRESH_HTTPONLY from environment to boolean: %v", err)
	}

	PasswordResetURLBase = getEnv("PASSWORD_RESET_URL_BASE", "http://localhost:3000/reset-password")
	PasswordResetTokenExpiresMinutes, err = strconv.Atoi(getEnv("PASSWORD_RESET_TOKEN_EXPIRES_MINUTES", "30"))
	if err != nil {
		log.Fatalf("unable to convert PASSWORD_RESET_TOKEN_EXPIRES_MINUTES from environment to integer: %v", err)
	}
	PasswordResetResendIntervalMinutes, err = strconv.Atoi(getEnv("PASSWORD_RESET_RESEND_INTERVAL_MINUTES", "5"))
	if err != nil {
		log.Fatalf("unable to convert PASSWORD_RESET_RESEND_INTERVAL_MINUTES from environment to integer: %v", err)
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
