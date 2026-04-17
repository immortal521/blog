package config

import (
	"fmt"
	"strings"
)

// validate performs basic validation on the loaded configuration.
//
// It ensures required fields are present and enforces environment-specific rules.
// In production mode, it additionally checks for sensitive configuration values
// such as database credentials and JWT secrets.
func validate(cfg *Config) error {
	var errs []string

	if cfg.App.Name == "" {
		errs = append(errs, "app.name is required")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		errs = append(errs, "server.port must be between 1 and 65535")
	}
	if cfg.App.IsProd() && cfg.JWT.Secret == "" {
		errs = append(errs, "jwt.secret is required in production (set JWT_SECRET env var)")
	}
	if cfg.App.IsProd() && cfg.Database.Password == "" {
		errs = append(errs, "database.password is required in production (set DATABASE_PASSWORD env var)")
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}
