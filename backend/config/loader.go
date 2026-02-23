package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	configInstance *Config
	configMutex    sync.RWMutex
	v              *viper.Viper
)

// Load initializes the configuration.
// It supports default values, file-based config, environment variables,
func Load(configPath string) (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	v = viper.New()
	setDefaults()

	// Enable environment variable overrides
	v.SetEnvPrefix("app")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	// Load primary configuration file
	if configPath != "" {
		if err := loadFromFile(configPath); err != nil {
			return nil, err
		}
	}

	// Handle environment-specific config files
	env := v.GetString("app.environment")
	if env != "" {
		ext := filepath.Ext(configPath)
		base := strings.TrimSuffix(configPath, ext)
		envFile := fmt.Sprintf("%s.%s%s", base, env, ext)
		if _, err := os.Stat(envFile); err == nil {
			v.SetConfigFile(envFile)
			if err := v.MergeInConfig(); err != nil {
				return nil, fmt.Errorf("failed to merge env config: %w", err)
			}
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.App.Domain = strings.TrimSuffix(cfg.App.Domain, "/")

	configInstance = cfg

	if configPath != "" {
		setupConfigWatch()
	}

	return cfg, nil
}

// setDefaults defines the fallback values for all configuration keys
func setDefaults() {
	v.SetDefault("app.name", "my-app")
	v.SetDefault("app.version", "1.0.1")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)
	v.SetDefault("app.domain", "localhost:8000")
	v.SetDefault("app.cors_origins", []string{"*"})

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8000)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.max_header_bytes", 1048576)
	v.SetDefault("server.graceful_shutdown", "10s")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.name", "blog")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "30m")
	v.SetDefault("database.timeout", "5s")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 2)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	v.SetDefault("jwt.secret", "your-default-jwt-secret-change-in-production")
	v.SetDefault("jwt.access_expiration", "15m")
	v.SetDefault("jwt.refresh_expiration", "168h")
	v.SetDefault("jwt.issuer", "my-app")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.compress", true)

	v.SetDefault("email.host", "mail.immort.top")
	v.SetDefault("email.port", 587)
	v.SetDefault("email.username", "Immortel@immort.top")
	v.SetDefault("email.password", "123456")
	v.SetDefault("email.from", "Immortel@immort.top")

	v.SetDefault("rustfs.region", "us-east-1")
	v.SetDefault("rustfs.access_key_id", "rustfsadmin")
	v.SetDefault("rustfs.secret_access_key", "rustfsadmin")
	v.SetDefault("rustfs.endpoint", "http://localhost:9000")
}

// loadFromFile reads the config file or creates a default one if it doesn't exist
func loadFromFile(configPath string) error {
	dir := filepath.Dir(configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0o755)
		defaultCfg := &Config{}
		if err := v.Unmarshal(defaultCfg); err != nil {
			return err
		}
		return writeDefaultConfig(configPath, defaultCfg)
	}

	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	return nil
}

// writeDefaultConfig serializes the current config state to a YAML file
func writeDefaultConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// setupConfigWatch enables hot-reloading when the configuration file is modified
func setupConfigWatch() {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)

		if err := v.ReadInConfig(); err != nil {
			log.Printf("Error reloading config file: %v", err)
			return
		}

		newCfg := &Config{}
		if err := v.Unmarshal(newCfg); err != nil {
			log.Printf("Error unmarshaling reloaded config: %v", err)
			return
		}

		if err := validateConfig(newCfg); err != nil {
			log.Printf("Reloaded config validation failed: %v", err)
			return
		}

		newCfg.App.Domain = strings.TrimSuffix(newCfg.App.Domain, "/")

		configMutex.Lock()
		configInstance = newCfg
		configMutex.Unlock()

		log.Println("Config reloaded successfully")
	})
}

func Get() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return configInstance
}

func GetViper() *viper.Viper {
	return v
}
