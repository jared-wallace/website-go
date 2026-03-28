package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration loaded from environment variables.
// All values come from the environment — no config files, no viper. (D-03)
type Config struct {
	DatabaseURL       string
	Port              string
	AppEnv            string
	AdminEmail        string // envOr("ADMIN_EMAIL", "")
	AdminPasswordHash string // envOr("ADMIN_PASSWORD_HASH", "")
	AdminHost         string // envOr("ADMIN_HOST", "admin.jared-wallace.com")
	SessionSecret     string // envOr("SESSION_SECRET", "") — reserved for future HMAC signing
	APIToken          string // envOr("API_TOKEN", "") — bearer token for push-to-publish API
	ImageDir          string // envOr("IMAGE_DIR", "/var/www/html/images") — EBS-backed image storage
}

// Load reads configuration from environment variables. Panics if required
// variables are missing (fail-fast is friendlier than a late nil-pointer).
func Load() Config {
	return Config{
		DatabaseURL:       mustEnv("DATABASE_URL"),
		Port:              envOr("PORT", "8080"),
		AppEnv:            envOr("APP_ENV", "development"),
		AdminEmail:        envOr("ADMIN_EMAIL", ""),
		AdminPasswordHash: envOr("ADMIN_PASSWORD_HASH", ""),
		AdminHost:         envOr("ADMIN_HOST", "admin.jared-wallace.com"),
		SessionSecret:     envOr("SESSION_SECRET", ""),
		APIToken:          envOr("API_TOKEN", ""),
		ImageDir:          envOr("IMAGE_DIR", "/var/www/html/images"),
	}
}

// envOr returns the value of the named environment variable, or fallback if
// the variable is unset or empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// mustEnv returns the value of the named environment variable. It panics if
// the variable is unset or empty — required config must be explicit.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}
