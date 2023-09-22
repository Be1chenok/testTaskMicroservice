package config

import "github.com/spf13/viper"

type Config struct {
	Server   serverConfig
	Database databaseConfig
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
	}, nil
}
