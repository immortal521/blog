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
	CorsOrigins []string `mapstructure:"corsorigins"`
}

type ServerConfig struct {
	Host             string        `mapstructure:"host"`
	Port             int           `mapstructure:"port"`
	ReadTimeout      time.Duration `mapstructure:"readtimeout"`
	WriteTimeout     time.Duration `mapstructure:"writetimeout"`
	IdleTimeout      time.Duration `mapstructure:"idletimeout"`
	MaxHeaderBytes   int           `mapstructure:"maxheaderbytes"`
	GracefulShutdown time.Duration `mapstructure:"gracefulshutdown"`
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
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"maxopenconns"`
	MaxIdleConns    int           `mapstructure:"maxidleconns"`
	ConnMaxLifetime time.Duration `mapstructure:"connmaxlifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connmaxidle_time"`
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
	PoolSize           int           `mapstructure:"poolsize"`
	MinIdleConns       int           `mapstructure:"minidleconns"`
	DialTimeout        time.Duration `mapstructure:"dialtimeout"`
	ReadTimeout        time.Duration `mapstructure:"readtimeout"`
	WriteTimeout       time.Duration `mapstructure:"writetimeout"`
	PoolTimeout        time.Duration `mapstructure:"pooltimeout"`
	IdleTimeout        time.Duration `mapstructure:"idletimeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idlecheckfrequency"`
}

type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	AccessExpiration  time.Duration `mapstructure:"accessexpiration"`
	RefreshExpiration time.Duration `mapstructure:"refreshexpiration"`
	Issuer            string        `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	FilePath   string `mapstructure:"filepath"`
	MaxSize    int    `mapstructure:"maxsize"`
	MaxBackups int    `mapstructure:"maxbackups"`
	MaxAge     int    `mapstructure:"maxage"`
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
	AccessKeyID     string `mapstructure:"accesskeyid"`
	SecretAccessKey string `mapstructure:"secretaccesskey"`
	Endpoint        string `mapstructure:"endpoint"`
}
