package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     uint   `mapstructure:"DB_PORT"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	DBHostReplica1     string `mapstructure:"DB_HOST_REPLICA_1"`
	DBPortReplica1     uint   `mapstructure:"DB_PORT_REPLICA_1"`
	DBUsernameReplica1 string `mapstructure:"DB_USERNAME_REPLICA_1"`
	DBPasswordReplica1 string `mapstructure:"DB_PASSWORD_REPLICA_1"`
	DBNameReplica1     string `mapstructure:"DB_NAME_REPLICA_1"`
	DBSSLModeReplica1  string `mapstructure:"DB_SSLMODE_REPLICA_1"`

	RedisHost    string `mapstructure:"REDIS_HOST"`
	RedisCacheDb int    `mapstructure:"REDIS_CACHE_DB"`

	RabbitMqAddr string `mapstructure:"RABBITMQ_ADDR"`

	FeedCache     time.Duration `mapstructure:"FEED_CACHE"`
	FeedLastCount uint          `mapstructure:"FEED_LAST_COUNT"`

	QueueWarmUpFeed string `mapstructure:"QUEUE_WARM_UP_FEED"`

	HttpServer string `mapstructure:"HTTP_SERVER"`
	Salt       string `mapstructure:"APP_SALT"`
	SigningKey string `mapstructure:"APP_SIGNING_KEY"`
	TokenTTL   uint   `mapstructure:"AUTH_TOKEN_TTL"`
}

func NewConfig() (*Config, error) {
	var config Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	logrus.Debug("Parsed config values")
	logrus.Debugln(config)

	return &config, nil
}
