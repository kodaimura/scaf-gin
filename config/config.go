package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	AppEnv  string
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
	LogLevel        string
	FrontendOrigins []string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("unable to load .env: %v", err)
	}

	AppEnv = requireChoice("APP_ENV", getEnv("APP_ENV", "dev"), "dev", "production", "test")
	AppName = getEnv("APP_NAME", "scaf-gin")
	AppHost = getEnv("APP_HOST", "localhost")
	AppPort = getEnv("APP_PORT", "8000")
	EnableSignup = parseBool("ENABLE_SIGNUP", getEnv("ENABLE_SIGNUP", "true"))
	AuthLoginIDMode = requireChoice("AUTH_LOGIN_ID_MODE", getEnv("AUTH_LOGIN_ID_MODE", "email"), "email", "login_id")

	DBEngine = requireChoice("DB_ENGINE", getEnv("DB_ENGINE", "postgres"), "postgres", "mysql", "sqlite3")
	DBName = getEnv("DB_NAME", "project_db")
	DBHost = getEnv("DB_HOST", "db")
	DBPort = getEnv("DB_PORT", "5432")
	DBUser = getEnv("DB_USER", "postgres")
	DBPass = getEnv("DB_PASSWORD", "postgres")

	MailProvider = requireChoice("MAIL_PROVIDER", getEnv("MAIL_PROVIDER", "mailhog"), "mailhog", "smtp")
	SMTPHost = getEnv("SMTP_HOST")
	SMTPPort = getEnv("SMTP_PORT", "587")
	SMTPUser = getEnv("SMTP_USER")
	SMTPPass = getEnv("SMTP_PASSWORD")
	MailFrom = getEnv("MAIL_FROM", "no-reply@example.local")

	BasicAuthUser = getEnv("BASIC_AUTH_USER")
	BasicAuthPass = getEnv("BASIC_AUTH_PASSWORD")

	AccessTokenSecret = getEnv("ACCESS_TOKEN_SECRET", "randomstring")
	AccessTokenExpiresSeconds = parseInt("ACCESS_TOKEN_EXPIRES_SECONDS", getEnv("ACCESS_TOKEN_EXPIRES_SECONDS", "900"))

	RefreshTokenSecret = getEnv("REFRESH_TOKEN_SECRET", "randomstring")
	RefreshTokenExpiresSeconds = parseInt("REFRESH_TOKEN_EXPIRES_SECONDS", getEnv("REFRESH_TOKEN_EXPIRES_SECONDS", "2592000"))
	RefreshTokenRememberMeExpiresSeconds = parseInt("REFRESH_TOKEN_REMEMBER_ME_EXPIRES_SECONDS", getEnv("REFRESH_TOKEN_REMEMBER_ME_EXPIRES_SECONDS", "2592000"))

	CookieAccessSecure = parseBool("COOKIE_ACCESS_SECURE", getEnv("COOKIE_ACCESS_SECURE", "true"))
	CookieRefreshSecure = parseBool("COOKIE_REFRESH_SECURE", getEnv("COOKIE_REFRESH_SECURE", "true"))
	CookieAccessHttpOnly = parseBool("COOKIE_ACCESS_HTTPONLY", getEnv("COOKIE_ACCESS_HTTPONLY", "true"))
	CookieRefreshHttpOnly = parseBool("COOKIE_REFRESH_HTTPONLY", getEnv("COOKIE_REFRESH_HTTPONLY", "true"))

	PasswordResetURLBase = getEnv("PASSWORD_RESET_URL_BASE", "http://localhost:3000/reset-password")
	PasswordResetTokenExpiresMinutes = parseInt("PASSWORD_RESET_TOKEN_EXPIRES_MINUTES", getEnv("PASSWORD_RESET_TOKEN_EXPIRES_MINUTES", "30"))
	PasswordResetResendIntervalMinutes = parseInt("PASSWORD_RESET_RESEND_INTERVAL_MINUTES", getEnv("PASSWORD_RESET_RESEND_INTERVAL_MINUTES", "5"))

	LogLevel = requireChoice("LOG_LEVEL", strings.ToLower(getEnv("LOG_LEVEL", "info")), "debug", "info", "warn", "error")
	FrontendOrigins = parseCSV(getEnv("FRONTEND_ORIGINS", "http://localhost:3000,http://localhost:5173"))

	validateConfig()
}

func validateConfig() {
	if len(FrontendOrigins) == 0 {
		log.Fatal("FRONTEND_ORIGINS must contain at least one origin")
	}
	if MailProvider == "smtp" && SMTPHost == "" {
		log.Fatal("SMTP_HOST is required when MAIL_PROVIDER=smtp")
	}
	if AppEnv != "production" {
		return
	}

	validateProductionSecret("ACCESS_TOKEN_SECRET", AccessTokenSecret)
	validateProductionSecret("REFRESH_TOKEN_SECRET", RefreshTokenSecret)
	for _, origin := range FrontendOrigins {
		if origin == "*" {
			log.Fatal("FRONTEND_ORIGINS must not contain * in production")
		}
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Scheme == "" || parsed.Hostname() == "" {
			log.Fatalf("invalid FRONTEND_ORIGINS value: %s", origin)
		}
		host := parsed.Hostname()
		if host == "localhost" || host == "127.0.0.1" || host == "::1" {
			log.Fatalf("FRONTEND_ORIGINS must not contain local origin in production: %s", origin)
		}
	}
}

func validateProductionSecret(key string, value string) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || normalized == "randomstring" || strings.HasPrefix(normalized, "change-me") {
		log.Fatalf("%s must be changed before production startup", key)
	}
}

func requireChoice(key string, value string, allowed ...string) string {
	normalized := strings.TrimSpace(value)
	for _, choice := range allowed {
		if normalized == choice {
			return normalized
		}
	}
	log.Fatalf("%s must be one of: %s", key, strings.Join(allowed, ", "))
	return ""
}

func parseBool(key string, value string) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("unable to convert %s from environment to boolean: %v", key, err)
	}
	return parsed
}

func parseInt(key string, value string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("unable to convert %s from environment to integer: %v", key, err)
	}
	return parsed
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
