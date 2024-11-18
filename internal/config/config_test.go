package config_test

import (
	"os"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/config"
)

func TestDBConfig_ConnectionString(t *testing.T) {
	cfg := config.DBConfig{
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		Host:     "localhost",
		Port:     5432,
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb"
	if got := cfg.ConnectionString(); got != expected {
		t.Errorf("ConnectionString() = %v, want %v", got, expected)
	}
}

func TestDBConfig_SafeConnectionString(t *testing.T) {
	cfg := config.DBConfig{
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		Host:     "localhost",
		Port:     5432,
	}

	expected := "postgres://testuser:****@localhost:5432/testdb"
	if got := cfg.SafeConnectionString(); got != expected {
		t.Errorf("SafeConnectionString() = %v, want %v", got, expected)
	}
}

func TestDBConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  config.DBConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: config.DBConfig{
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				Host:     "localhost",
				Port:     5432,
			},
			wantErr: false,
		},
		{
			name: "missing password",
			config: config.DBConfig{
				User:   "testuser",
				DBName: "testdb",
				Host:   "localhost",
				Port:   5432,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too low",
			config: config.DBConfig{
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				Host:     "localhost",
				Port:     0,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: config.DBConfig{
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				Host:     "localhost",
				Port:     65536,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Setup test environment variables
	envVars := map[string]string{
		"DB_USER":     "testuser",
		"DB_PASSWORD": "testpass",
		"DB_NAME":     "testdb",
		"DB_HOST":     "testhost",
		"DB_PORT":     "5432",
	}

	// Set environment variables
	for k, v := range envVars {
		os.Setenv(k, v)
	}

	// Cleanup environment variables after test
	defer func() {
		for k := range envVars {
			os.Unsetenv(k)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DB.User != envVars["DB_USER"] {
		t.Errorf("Load() User = %v, want %v", cfg.DB.User, envVars["DB_USER"])
	}
	if cfg.DB.Password != envVars["DB_PASSWORD"] {
		t.Errorf("Load() Password = %v, want %v", cfg.DB.Password, envVars["DB_PASSWORD"])
	}
	if cfg.DB.DBName != envVars["DB_NAME"] {
		t.Errorf("Load() DBName = %v, want %v", cfg.DB.DBName, envVars["DB_NAME"])
	}
	if cfg.DB.Host != envVars["DB_HOST"] {
		t.Errorf("Load() Host = %v, want %v", cfg.DB.Host, envVars["DB_HOST"])
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	if got := config.GetEnvOrDefault("TEST_KEY", "default"); got != "test_value" {
		t.Errorf("getEnvOrDefault() = %v, want %v", got, "test_value")
	}

	// Test with environment variable not set
	if got := config.GetEnvOrDefault("NONEXISTENT_KEY", "default"); got != "default" {
		t.Errorf("getEnvOrDefault() = %v, want %v", got, "default")
	}
} 