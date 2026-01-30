package config

import (
	"fmt"
	"time"
)

func validateConfig(cfg *Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if cfg.JWT.Secret == "your-default-jwt-secret-change-in-production" &&
		cfg.App.Environment == "production" {
		return fmt.Errorf("JWT secret must be changed in production environment")
	}

	if cfg.JWT.AccessExpiration < time.Minute {
		return fmt.Errorf("JWT access expiration must be at least 1 minute")
	}

	if cfg.Redis.PoolSize <= 0 {
		return fmt.Errorf("redis pool size must be positive")
	}

	return nil
}
