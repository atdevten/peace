package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the overall configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	App      AppConfig      `yaml:"app"`
	Auth     AuthConfig     `yaml:"auth"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string        `yaml:"port" env:"SERVER_PORT"`
	Host         string        `yaml:"host" env:"SERVER_HOST"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

// PostgresConfig represents PostgreSQL configuration
type PostgresConfig struct {
	Host            string        `yaml:"host" env:"POSTGRES_HOST"`
	Port            string        `yaml:"port" env:"POSTGRES_PORT"`
	User            string        `yaml:"user" env:"POSTGRES_USER"`
	Password        string        `yaml:"password" env:"POSTGRES_PASSWORD"`
	DBName          string        `yaml:"dbname" env:"POSTGRES_DB"`
	SSLMode         string        `yaml:"sslmode" env:"POSTGRES_SSLMODE"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     string `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

// AppConfig represents application configuration
type AppConfig struct {
	Environment string     `yaml:"environment" env:"ENVIRONMENT"`
	LogLevel    string     `yaml:"log_level" env:"LOG_LEVEL"`
	CORS        CORSConfig `yaml:"cors"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWT    JWTConfig    `yaml:"jwt"`
	Google GoogleConfig `yaml:"google"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret            string        `yaml:"secret" env:"JWT_SECRET"`
	Expiration        time.Duration `yaml:"expiration" env:"JWT_EXPIRATION"`
	RefreshExpiration time.Duration `yaml:"refresh_expiration" env:"JWT_REFRESH_EXPIRATION"`
}

// GoogleConfig represents Google OAuth configuration
type GoogleConfig struct {
	ClientID     string `yaml:"client_id" env:"GOOGLE_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" env:"GOOGLE_CLIENT_SECRET"`
	RedirectURI  string `yaml:"redirect_uri" env:"GOOGLE_REDIRECT_URI"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string `yaml:"level" env:"LOG_LEVEL"`
	Format     string `yaml:"format" env:"LOG_FORMAT"` // "json" or "pretty"
	TimeFormat string `yaml:"time_format" env:"LOG_TIME_FORMAT"`
	Caller     bool   `yaml:"caller" env:"LOG_CALLER"`
	CallerSkip int    `yaml:"caller_skip" env:"LOG_CALLER_SKIP"`
}

// Load loads configuration from YAML file and environment variables
func Load() (*Config, error) {
	return LoadWithPath("")
}

// LoadWithPath loads configuration from a specific YAML file path and environment variables
func LoadWithPath(configPath string) (*Config, error) {
	config := &Config{}

	// Load YAML config
	if err := loadYAMLConfigWithPath(config, configPath); err != nil {
		return nil, fmt.Errorf("failed to load YAML config: %w", err)
	}

	return config, nil
}

// loadYAMLConfigWithPath loads configuration from a specific YAML file path
func loadYAMLConfigWithPath(config *Config, configPath string) error {
	if configPath == "" {
		return fmt.Errorf("config path is empty")
	}

	if configData, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(configData, config); err != nil {
			return fmt.Errorf("failed to unmarshal YAML from %s: %w", configPath, err)
		}
		return nil
	}

	return nil
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Postgres.Host,
		c.Database.Postgres.Port,
		c.Database.Postgres.User,
		c.Database.Postgres.Password,
		c.Database.Postgres.DBName,
		c.Database.Postgres.SSLMode,
	)
}

// GetRedisAddr returns the Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Database.Redis.Host, c.Database.Redis.Port)
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetLogLevel returns the configured log level as zerolog.Level
func (c *Config) GetLogLevel() string {
	return c.Log.Level
}

// GetLogFormat returns the configured log format
func (c *Config) GetLogFormat() string {
	return c.Log.Format
}

// IsPrettyLogging returns true if logging should be in pretty format
func (c *Config) IsPrettyLogging() bool {
	return c.Log.Format == "pretty"
}

// GetLogTimeFormat returns the configured time format for logs
func (c *Config) GetLogTimeFormat() string {
	return c.Log.TimeFormat
}

// ShouldLogCaller returns true if caller information should be included in logs
func (c *Config) ShouldLogCaller() bool {
	return c.Log.Caller
}

// GetLogCallerSkip returns the number of stack frames to skip when logging caller
func (c *Config) GetLogCallerSkip() int {
	return c.Log.CallerSkip
}
