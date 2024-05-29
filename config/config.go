package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	DataFile       string
	WindowDuration time.Duration
	PersistInterval time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{
		Port:           viper.GetString("port"),
		DataFile:       viper.GetString("data_file"),
		WindowDuration: viper.GetDuration("window_duration"),
		PersistInterval: viper.GetDuration("persist_interval"),
	}

	return cfg, nil
}