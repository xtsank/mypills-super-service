package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"log/slog"

	"github.com/joho/godotenv"
	"github.com/samber/do/v2"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBHost     string

	JWTSecret     string
	TokenDuration time.Duration

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	ServerAddress string

	LogFile  string
	LogLevel slog.Level
	logger   *slog.Logger
}

func NewConfig(i do.Injector) (*Config, error) {
	logger := do.MustInvoke[*slog.Logger](i)

	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, reading from system environment")
	}

	durationInt := parseIntEnv(logger, "TOKEN_DURATION", 86400)

	maxOpen := parseIntEnv(logger, "DB_MAX_OPEN_CONNS", 25)
	maxIdle := parseIntEnv(logger, "DB_MAX_IDLE_CONNS", 10)
	maxLifetime := parseDurationEnv(logger, "DB_CONN_MAX_LIFETIME", time.Hour)
	maxIdleTime := parseDurationEnv(logger, "DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	logLevel := parseLogLevelEnv(logger, "LOG_LEVEL", slog.LevelDebug)
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "logs/app.log"
	}

	return &Config{
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBPort:        os.Getenv("DB_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		JWTSecret:     os.Getenv("SECRET_KEY"),
		TokenDuration: time.Duration(durationInt) * time.Second,

		DBMaxOpenConns:    maxOpen,
		DBMaxIdleConns:    maxIdle,
		DBConnMaxLifetime: maxLifetime,
		DBConnMaxIdleTime: maxIdleTime,

		ServerAddress: os.Getenv("SERVER_PORT"),

		LogFile:  logFile,
		LogLevel: logLevel,
		logger:   logger,
	}, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) GetJWTConfig() (string, time.Duration) {
	if c.JWTSecret == "" {
		c.logger.Error("JWT_SECRET is not set in environment")
		os.Exit(1)
	}

	return c.JWTSecret, c.TokenDuration
}

func parseIntEnv(logger *slog.Logger, key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		logger.Warn("Invalid env, using default", slog.String("key", key), slog.Any("error", err), slog.Int("default", fallback))
		return fallback
	}
	return parsed
}

func parseDurationEnv(logger *slog.Logger, key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		logger.Warn("Invalid env, using default", slog.String("key", key), slog.Any("error", err), slog.String("default", fallback.String()))
		return fallback
	}
	return parsed
}

func parseLogLevelEnv(logger *slog.Logger, key string, fallback slog.Level) slog.Level {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if val == "" {
		return fallback
	}

	switch val {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		logger.Warn("Invalid log level, using default", slog.String("key", key), slog.String("value", val), slog.String("default", fallback.String()))
		return fallback
	}
}
