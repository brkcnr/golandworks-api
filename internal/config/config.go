package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/brkcnr/golandworks-api/internal/apierror"
	"github.com/joho/godotenv"
)

// DBConfig is the database configuration.
type DBConfig struct {
	User string

	Password string

	DBName string

	Host string

	Port int
}

// Config is the application configuration.
type Config struct {
	DB DBConfig
}

// ConnectionString returns the full database connection string.
func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf(

		"postgres://%s:%s@%s:%d/%s",

		c.User,

		c.Password,

		c.Host,

		c.Port,

		c.DBName,
	)
}

// SafeConnectionString returns a connection string with sensitive data redacted.
func (c DBConfig) SafeConnectionString() string {
	return fmt.Sprintf(

		"postgres://%s:****@%s:%d/%s",

		c.User,

		c.Host,

		c.Port,

		c.DBName,
	)
}

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: *dbConfig,
	}, nil
}

// loadDBConfig loads the database configuration from environment variables.
func loadDBConfig() (*DBConfig, error) {
	port, portErr := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))

	if portErr != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", portErr)
	}
	config := &DBConfig{
		User: getEnvOrDefault("DB_USER", "postgres"),

		Password: os.Getenv("DB_PASSWORD"),

		DBName: getEnvOrDefault("DB_NAME", "postgres"),

		Host: getEnvOrDefault("DB_HOST", "localhost"),

		Port: port,
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate validates the database configuration.
func (c DBConfig) validate() error {
	if c.Password == "" {
		return apierror.Wrap(

			apierror.ErrMissingDBPassword,

			http.StatusBadRequest,

			"database password is required",
		)
	}

	if c.Port < 1 || c.Port > 65535 {
		return apierror.Wrap(

			apierror.ErrInvalidDBPort,

			http.StatusBadRequest,

			fmt.Sprintf("port value %d is not between 1 and 65535", c.Port),
		)
	}

	return nil
}

// getEnvOrDefault gets the environment variable or the default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
