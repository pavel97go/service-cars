package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App struct {
		Port string
	}
	Metrics struct {
		Port string
	}
	DB struct {
		Port     string
		Host     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Cache struct {
		TTLSeconds int
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
func Init() *Config {
	var c Config
	c.App.Port = env("APP_PORT", "8080")

	c.DB.Host = env("DB_HOST", "localhost")
	c.DB.Port = env("DB_PORT", "5432")
	c.DB.User = env("DB_USER", "cars")
	c.DB.Password = env("DB_PASSWORD", "cars")
	c.DB.Name = env("DB_NAME", "cars")
	c.DB.SSLMode = env("DB_SSLMODE", "disable")

	c.Metrics.Port = env("METRICS_PORT", "9100")
	c.Cache.TTLSeconds = envInt("CACHE_TTL_SECONDS", 60)

	return &c
}
func (c *Config) GetConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode,
	)
}
