package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Worker   WorkerConfig
}

type ServerConfig struct {
	Port             int
	ReadTimeoutSecs  int
	WriteTimeoutSecs int
	IdleTimeoutSecs  int
	AllowedOrigins   []string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type WorkerConfig struct {
	Concurrency  int
	PollInterval int
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:             getEnvInt("PORT"),
			ReadTimeoutSecs:  getEnvInt("READ_TIMEOUT_SECONDS"),
			WriteTimeoutSecs: getEnvInt("WRITE_TIMEOUT_SECONDS"),
			IdleTimeoutSecs:  getEnvInt("IDLE_TIMEOUT_SECONDS"),
			AllowedOrigins:   getEnvList("CORS_ALLOWED_ORIGINS"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnvInt("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR"),
			Password: getEnv("REDIS_PASSWORD"),
			DB:       getEnvInt("REDIS_DB"),
		},
		Worker: WorkerConfig{
			Concurrency:  getEnvInt("WORKER_CONCURRENCY"),
			PollInterval: getEnvInt("WORKER_POLL_INTERVAL"),
		},
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvInt(key string) int {
	value, _ := strconv.Atoi(os.Getenv(key))
	return value
}

func getEnvList(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}
