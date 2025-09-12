package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the overall configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
	Auth     AuthConfig
	Log      LogConfig
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port          string
	Host          string
	WebSocketPort string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

// PostgresConfig represents PostgreSQL configuration
type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// AppConfig represents application configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	CORS        CORSConfig
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWT    JWTConfig
	Google GoogleConfig
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

// GoogleConfig represents Google OAuth configuration
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string // "json" or "pretty"
	Format     string
	TimeFormat string
	Caller     bool
	CallerSkip int
}

// Load loads configuration from environment variables and optional .env file
func Load() (*Config, error) {
	return LoadWithEnvFile("")
}

// LoadWithEnvFile loads configuration from .env file and environment variables
func LoadWithEnvFile(envFilePath string) (*Config, error) {
	// Load .env file if provided (optional)
	if envFilePath != "" {
		if err := godotenv.Load(envFilePath); err != nil {
			// Don't fail if .env file doesn't exist, just continue with system env vars
			fmt.Printf("Warning: Could not load .env file %s: %v\n", envFilePath, err)
		}
	}

	// Load configuration from environment variables
	config, err := loadFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	return config, nil
}

// LoadWithPath loads configuration from .env file (backward compatibility)
func LoadWithPath(configPath string) (*Config, error) {
	// Convert YAML path to .env path for backward compatibility
	envPath := strings.Replace(configPath, ".yml", ".env", 1)
	return LoadWithEnvFile(envPath)
}

// loadFromEnvironment loads configuration from environment variables
func loadFromEnvironment() (*Config, error) {
	config := &Config{}

	// Load server config
	config.Server.Port = getEnvOrDefault("SERVER_PORT", "8080")
	config.Server.Host = getEnvOrDefault("SERVER_HOST", "0.0.0.0")
	config.Server.WebSocketPort = getEnvOrDefault("WEBSOCKET_PORT", "8081")

	var err error
	config.Server.ReadTimeout, err = time.ParseDuration(getEnvOrDefault("SERVER_READ_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
	}
	config.Server.WriteTimeout, err = time.ParseDuration(getEnvOrDefault("SERVER_WRITE_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
	}
	config.Server.IdleTimeout, err = time.ParseDuration(getEnvOrDefault("SERVER_IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_IDLE_TIMEOUT: %w", err)
	}

	// Load database config - Postgres
	config.Database.Postgres.Host = getEnvOrDefault("POSTGRES_HOST", "localhost")
	config.Database.Postgres.Port = getEnvOrDefault("POSTGRES_PORT", "5432")
	config.Database.Postgres.User = getEnvOrDefault("POSTGRES_USER", "postgres")
	config.Database.Postgres.Password = getEnvOrDefault("POSTGRES_PASSWORD", "")
	config.Database.Postgres.DBName = getEnvOrDefault("POSTGRES_DB", "peace")
	config.Database.Postgres.SSLMode = getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	config.Database.Postgres.MaxOpenConns = getEnvAsIntOrDefault("POSTGRES_MAX_OPEN_CONNS", 25)
	config.Database.Postgres.MaxIdleConns = getEnvAsIntOrDefault("POSTGRES_MAX_IDLE_CONNS", 5)
	config.Database.Postgres.ConnMaxLifetime, err = time.ParseDuration(getEnvOrDefault("POSTGRES_CONN_MAX_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_CONN_MAX_LIFETIME: %w", err)
	}

	// Load database config - Redis
	config.Database.Redis.Host = getEnvOrDefault("REDIS_HOST", "localhost")
	config.Database.Redis.Port = getEnvOrDefault("REDIS_PORT", "6379")
	config.Database.Redis.Password = getEnvOrDefault("REDIS_PASSWORD", "")
	config.Database.Redis.DB = getEnvAsIntOrDefault("REDIS_DB", 0)

	// Load app config
	config.App.Environment = getEnvOrDefault("ENVIRONMENT", "development")
	config.App.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")

	// Parse CORS origins
	corsOrigins := getEnvOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost")
	config.App.CORS.AllowedOrigins = strings.Split(corsOrigins, ",")
	for i, origin := range config.App.CORS.AllowedOrigins {
		config.App.CORS.AllowedOrigins[i] = strings.TrimSpace(origin)
	}

	// Parse CORS methods
	corsMethods := getEnvOrDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	config.App.CORS.AllowedMethods = strings.Split(corsMethods, ",")
	for i, method := range config.App.CORS.AllowedMethods {
		config.App.CORS.AllowedMethods[i] = strings.TrimSpace(method)
	}

	// Load JWT config
	config.Auth.JWT.Secret = getEnvOrDefault("JWT_SECRET", "dev-secret-key")
	config.Auth.JWT.Expiration, err = time.ParseDuration(getEnvOrDefault("JWT_EXPIRATION", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}
	config.Auth.JWT.RefreshExpiration, err = time.ParseDuration(getEnvOrDefault("JWT_REFRESH_EXPIRATION", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRATION: %w", err)
	}

	// Load Google OAuth config
	config.Auth.Google.ClientID = getEnvOrDefault("GOOGLE_CLIENT_ID", "")
	config.Auth.Google.ClientSecret = getEnvOrDefault("GOOGLE_CLIENT_SECRET", "")
	config.Auth.Google.RedirectURI = getEnvOrDefault("GOOGLE_REDIRECT_URI", "http://localhost:3000/auth/google/callback")

	// Load log config
	config.Log.Level = getEnvOrDefault("LOG_LEVEL", "info")
	config.Log.Format = getEnvOrDefault("LOG_FORMAT", "pretty")
	config.Log.TimeFormat = getEnvOrDefault("LOG_TIME_FORMAT", "2006-01-02T15:04:05Z07:00")
	config.Log.Caller = getEnvAsBoolOrDefault("LOG_CALLER", false)
	config.Log.CallerSkip = getEnvAsIntOrDefault("LOG_CALLER_SKIP", 2)

	return config, nil
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
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
