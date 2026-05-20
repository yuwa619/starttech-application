package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment   string
	ServiceName   string
	Port          string
	MongoDBURI    string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisTLS      bool
}

func Load() Config {
	return Config{
		Environment:   getEnv("APP_ENV", "local"),
		ServiceName:   getEnv("SERVICE_NAME", "starttech-api"),
		Port:          getEnv("PORT", "8080"),
		MongoDBURI:    os.Getenv("MONGODB_URI"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		RedisTLS:      getEnvBool("REDIS_TLS", false),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
