package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server       serverConfig
	Database     databaseConfig
	Cache        cacheConfig
	UserPassword userPasswordConfig
	Tokens       tokensConfig
}

type serverConfig struct {
	Host string
	Port string
}

type databaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type cacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type userPasswordConfig struct {
	Salt string
}

type tokensConfig struct {
	SigningKey string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Init() (*Config, error) {
	viper.SetConfigFile("./../../.env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		serverConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetString("SERVER_PORT"),
		},
		databaseConfig{
			Host:     viper.GetString("PG_HOST"),
			Port:     viper.GetString("PG_PORT"),
			Username: viper.GetString("PG_USER"),
			Password: viper.GetString("PG_PASSWORD"),
			DBName:   viper.GetString("PG_BASE"),
			SSLMode:  viper.GetString("PG_SSL_MODE"),
		},
		cacheConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		userPasswordConfig{
			Salt: viper.GetString("USER_PASSWORD_SALT"),
		},
		tokensConfig{
			SigningKey: viper.GetString("TOKENS_SIGNING_KEY"),
			AccessTTL:  viper.GetDuration("ACCESS_TOKEN_TTL") * time.Second,
			RefreshTTL: viper.GetDuration("REFRESH_TOKEN_TTL") * time.Second,
		},
	}, nil
}
