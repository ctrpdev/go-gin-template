package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBSSLMode     string `mapstructure:"DB_SSLMODE"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Error reading config file. Fallback to environment variables", "err", err)
	}

	// 3. Bind explicitly environment variables keys
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("JWT_SECRET_KEY")

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

// GetDatabaseURL constructs the PostgreSQL connection string
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

// GetRedisURL constructs the Redis connection string
func (c *Config) GetRedisURL() string {
	return c.RedisHost + ":" + c.RedisPort
}
