package config

import (
	"os"
	"testing"
)

func TestLoad_MissingDatabaseURL_Panics(t *testing.T) {
	t.Helper()
	// Ensure DATABASE_URL is unset for this test.
	original := os.Getenv("DATABASE_URL")
	os.Unsetenv("DATABASE_URL")
	defer func() {
		if original != "" {
			os.Setenv("DATABASE_URL", original)
		}
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Load() to panic when DATABASE_URL is unset, but it did not")
		}
	}()

	Load()
}

func TestLoad_Defaults(t *testing.T) {
	// Set only the required variable; leave PORT and APP_ENV unset.
	os.Setenv("DATABASE_URL", "postgres://test/testdb")
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")
	defer func() {
		os.Unsetenv("DATABASE_URL")
	}()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default Port %q, got %q", "8080", cfg.Port)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("expected default AppEnv %q, got %q", "development", cfg.AppEnv)
	}
}

func TestEnvOr_Override(t *testing.T) {
	os.Setenv("TEST_ENVVAR_OVERRIDE", "custom-value")
	defer os.Unsetenv("TEST_ENVVAR_OVERRIDE")

	got := envOr("TEST_ENVVAR_OVERRIDE", "default-value")
	if got != "custom-value" {
		t.Errorf("expected %q, got %q", "custom-value", got)
	}
}
