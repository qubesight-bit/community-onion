package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host         string
	Port         string
	Name         string
	User         string
	PasswordFile string
}

func ConfigFromEnvironment() Config {
	return Config{
		Host:         envOrDefault("DB_HOST", "database"),
		Port:         envOrDefault("DB_PORT", "5432"),
		Name:         envOrDefault("DB_NAME", "community_onion"),
		User:         envOrDefault("DB_USER", "community_app"),
		PasswordFile: envOrDefault("DB_PASSWORD_FILE", "/run/secrets/postgres_app_password"),
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	passwordBytes, err := os.ReadFile(cfg.PasswordFile)
	if err != nil {
		return nil, fmt.Errorf("read database password: %w", err)
	}

	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return nil, fmt.Errorf("database password is empty")
	}

	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	query := connectionURL.Query()
	query.Set("sslmode", "disable")
	connectionURL.RawQuery = query.Encode()

	poolConfig, err := pgxpool.ParseConfig(connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
