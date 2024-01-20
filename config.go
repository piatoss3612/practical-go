package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv           string
	RedisURL        string
	DockerRedisURL  string
	DiscordBotToken string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		GoEnv:           getEnv("GO_ENV", "development"),
		RedisURL:        getEnv("REDIS_URL", "localhost:6379"),
		DockerRedisURL:  getEnv("DOCKER_REDIS_URL", "redis:6379"),
		DiscordBotToken: getEnv("DISCORD_BOT_TOKEN", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
