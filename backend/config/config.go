// Package config provides application configuration structures and helpers
// for loading and accessing runtime settings.
package config

import (
	"fmt"
	"time"
)

const (
	EnvProd = "production"
	EnvDev  = "development"
)

// Config represents the root configuration structure of the application.
// It aggregates all subsystem configurations.
type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app" toml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server" toml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database" toml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis" toml:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt" yaml:"jwt" toml:"jwt"`
	Log      LogConfig      `mapstructure:"log" yaml:"log" toml:"log"`
	Email    EmailConfig    `mapstructure:"email" yaml:"email" toml:"email"`
	LLM      LLMConfig      `mapstructure:"llm" yaml:"llm" toml:"llm"`
	Rustfs   RustfsConfig   `mapstructure:"rustfs" yaml:"rustfs" toml:"rustfs"`
}

// AppConfig contains general application-level settings such as environment,
// version, and runtime behavior flags.
type AppConfig struct {
	Name        string   `mapstructure:"name" yaml:"name" toml:"name"`
	Version     string   `mapstructure:"version" yaml:"version" toml:"version"`
	Environment string   `mapstructure:"environment" yaml:"environment" toml:"environment"`
	Debug       bool     `mapstructure:"debug" yaml:"debug" toml:"debug"`
	Domain      string   `mapstructure:"domain" yaml:"domain" toml:"domain"`
	CorsOrigins []string `mapstructure:"cors_origins" yaml:"cors_origins" toml:"cors_origins"`
}

// IsProd reports whether the application is running in production environment.
func (a AppConfig) IsProd() bool {
	return a.Environment == EnvProd
}

// IsDev reports whether the application is running in development environment.
func (a AppConfig) IsDev() bool {
	return a.Environment == EnvDev
}

// ServerConfig defines HTTP server configuration including network settings
// and timeout controls.
type ServerConfig struct {
	Host             string        `mapstructure:"host" yaml:"host" toml:"host"`
	Port             int           `mapstructure:"port" yaml:"port" toml:"port"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout"`
	MaxHeaderBytes   int           `mapstructure:"max_header_bytes" yaml:"max_header_bytes" toml:"max_header_bytes"`
	GracefulShutdown time.Duration `mapstructure:"graceful_shutdown" yaml:"graceful_shutdown" toml:"graceful_shutdown"`
}

// Addr returns the server address in "host:port" format.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// DatabaseConfig contains relational database connection and pool settings.
type DatabaseConfig struct {
	Host            string        `mapstructure:"host" yaml:"host" toml:"host"`
	Port            int           `mapstructure:"port" yaml:"port" toml:"port"`
	User            string        `mapstructure:"user" yaml:"user" toml:"user"`
	Password        string        `mapstructure:"password" yaml:"password" toml:"password"`
	Name            string        `mapstructure:"name" yaml:"name" toml:"name"`
	SSLMode         string        `mapstructure:"ssl_mode" yaml:"ssl_mode" toml:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns" toml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns" toml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" toml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time" toml:"conn_max_idle_time"`
	Timeout         time.Duration `mapstructure:"timeout" yaml:"timeout" toml:"timeout"`
}

// RedisConfig defines Redis connection and pool configuration.
type RedisConfig struct {
	Host               string        `mapstructure:"host" yaml:"host" toml:"host"`
	Port               int           `mapstructure:"port" yaml:"port" toml:"port"`
	Password           string        `mapstructure:"password" yaml:"password" toml:"password"`
	DB                 int           `mapstructure:"db" yaml:"db" toml:"db"`
	PoolSize           int           `mapstructure:"pool_size" yaml:"pool_size" toml:"pool_size"`
	MinIdleConns       int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns" toml:"min_idle_conns"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	PoolTimeout        time.Duration `mapstructure:"pool_timeout" yaml:"pool_timeout" toml:"pool_timeout"`
	IdleTimeout        time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idle_check_frequency" yaml:"idle_check_frequency" toml:"idle_check_frequency"`
}

// JWTConfig contains configuration for JWT authentication and token lifecycle.
type JWTConfig struct {
	Secret            string        `mapstructure:"secret" yaml:"secret" toml:"secret"`
	AccessExpiration  time.Duration `mapstructure:"access_expiration" yaml:"access_expiration" toml:"access_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration" yaml:"refresh_expiration" toml:"refresh_expiration"`
	Issuer            string        `mapstructure:"issuer" yaml:"issuer" toml:"issuer"`
}

// LogConfig defines logging configuration including output format and rotation policy.
type LogConfig struct {
	Level      string `mapstructure:"level" yaml:"level" toml:"level"`
	Format     string `mapstructure:"format" yaml:"format" toml:"format"`
	FilePath   string `mapstructure:"file_path" yaml:"file_path" toml:"file_path"`
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size" toml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups" toml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age" toml:"max_age"`
	Compress   bool   `mapstructure:"compress" yaml:"compress" toml:"compress"`
}

// EmailConfig contains SMTP configuration for sending emails.
type EmailConfig struct {
	Host     string `mapstructure:"host" yaml:"host" toml:"host"`
	Port     int    `mapstructure:"port" yaml:"port" toml:"port"`
	Username string `mapstructure:"username" yaml:"username" toml:"username"`
	Password string `mapstructure:"password" yaml:"password" toml:"password"`
	From     string `mapstructure:"from" yaml:"from" toml:"from"`
}

// LLMConfig contains configuration for large language model integration.
type LLMConfig struct {
	APIKey string `mapstructure:"apikey" yaml:"apikey" toml:"apikey"`
}

// RustfsConfig contains configuration for RustFS or S3-compatible object storage service.
type RustfsConfig struct {
	Region          string `mapstructure:"region" yaml:"region" toml:"region"`
	AccessKeyID     string `mapstructure:"access_key_id" yaml:"access_key_id" toml:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" yaml:"secret_access_key" toml:"secret_access_key"`
	Endpoint        string `mapstructure:"endpoint" yaml:"endpoint" toml:"endpoint"`
}
