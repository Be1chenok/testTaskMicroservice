package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Host string
	Port string
}

func SrvInit() (*ServerConfig, error) {
	viper.SetConfigFile("./../../.env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &ServerConfig{
		Host: viper.GetString("SERVER_HOST"),
		Port: viper.GetString("SERVER_PORT"),
	}, nil
}
