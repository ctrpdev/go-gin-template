package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	JWTSecret     string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Error reading config file. Fallback to environment variables", "err", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Default values
	if config.ServerAddress == "" {
		config.ServerAddress = ":8080"
	}

	return &config, nil
}
