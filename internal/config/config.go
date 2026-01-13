package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	JWT JWTConfig
}

type ServerConfig struct {
	Host string
	Port string
	Env string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
	Secret string
	ExpiryHour int
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
			Env: getEnv("APP_ENV", "dev"),
		},
		Database: DatabaseConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName: getEnv("DB_NAME", ""),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			ExpiryHour: getEnvAsInt("JWT_EXPIRY_HOUR", 24),
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
    if c.Database.Password == "" {
        return fmt.Errorf("DB_PASSWORD is required")
    }
    if c.JWT.Secret == "" {
        return fmt.Errorf("JWT_SECRET is required")
    }
    return nil
}

func (d DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s?sslmode=%s",
        d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
    )
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultVal int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultVal
}

