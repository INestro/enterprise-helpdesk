package config

import "os"

type Config struct {
	HTTPAddr   string
	PostgreDSN string
	RedisAddr  string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		HTTPAddr:   getEnv("HTTP_ADDR", ":8080"),
		PostgreDSN: getEnv("POSTGRE_DSN", "postgres://user:pass@localhost:5432/enterprise-helpdesk?sslmode=disable"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
