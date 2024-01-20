package main

import "github.com/spf13/viper"

type Config struct {
	RedisURL        string `mapstructure:"REDIS_URL"`
	DockerRedisURL  string `mapstructure:"DOCKER_REDIS_URL"`
	GoEnv           string `mapstructure:"GO_ENV"`
	DiscordBotToken string `mapstructure:"DISCORD_BOT_TOKEN"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
