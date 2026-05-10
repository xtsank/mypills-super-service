package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
}

func NewConfig(i do.Injector) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	durationStr := os.Getenv("TOKEN_DURATION")
	durationInt, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Printf("Invalid TOKEN_DURATION: %v. Using default 24h", err)
		durationInt = 86400
	}

	maxOpen := parseIntEnv("DB_MAX_OPEN_CONNS", 25)
	maxIdle := parseIntEnv("DB_MAX_IDLE_CONNS", 10)
	maxLifetime := parseDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour)
	maxIdleTime := parseDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

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
	}, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) GetJWTConfig() (string, time.Duration) {
	if c.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}

	return c.JWTSecret, c.TokenDuration
}

func parseIntEnv(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Invalid %s: %v. Using default %d", key, err, fallback)
		return fallback
	}
	return parsed
}

func parseDurationEnv(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("Invalid %s: %v. Using default %s", key, err, fallback)
		return fallback
	}
	return parsed
}
