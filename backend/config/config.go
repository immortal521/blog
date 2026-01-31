// Package config provider config
package config

import (
	"fmt"
	"time"
)

const (
	EnvProd = "production"
	EnvDev  = "development"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Email    EmailConfig    `mapstructure:"email"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Rustfs   RustfsConfig   `mapstructure:"rustfs"`
}

type AppConfig struct {
	Name        string   `mapstructure:"name"`
	Version     string   `mapstructure:"version"`
	Environment string   `mapstructure:"environment"`
	Debug       bool     `mapstructure:"debug"`
	Domain      string   `mapstructure:"domain"`
	CorsOrigins []string `mapstructure:"cors_origins"`
}

type ServerConfig struct {
	Host             string        `mapstructure:"host"`
	Port             int           `mapstructure:"port"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes   int           `mapstructure:"max_header_bytes"`
	GracefulShutdown time.Duration `mapstructure:"graceful_shutdown"`
}

func (s ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	Timeout         time.Duration `mapstructure:"timeout"`
}

func (d DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable connect_timeout=%d",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		int(d.Timeout.Seconds()),
	)
}

type RedisConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	Password           string        `mapstructure:"password"`
	DB                 int           `mapstructure:"db"`
	PoolSize           int           `mapstructure:"pool_size"`
	MinIdleConns       int           `mapstructure:"min_idle_conns"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
	PoolTimeout        time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idle_check_frequency"`
}

type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	AccessExpiration  time.Duration `mapstructure:"access_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
	Issuer            string        `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type LLMConfig struct {
	APIKey string `mapstructure:"apikey"`
}

type RustfsConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Endpoint        string `mapstructure:"endpoint"`
}
