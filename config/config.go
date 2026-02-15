package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerURL       string
	Password        string
	PollIntervalSec int
	MessageLimit    int
	ChatLimit       int
}

func Load() (*Config, error) {
	viper.SetConfigName("bluebubbles")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/bluebubbles-tui/")
	viper.AddConfigPath(".")

	// Env var bindings
	viper.SetEnvPrefix("BB")
	viper.AutomaticEnv()
	viper.BindEnv("server_url", "BB_SERVER_URL")
	viper.BindEnv("password", "BB_PASSWORD")

	// Defaults
	viper.SetDefault("poll_interval_sec", 10)
	viper.SetDefault("message_limit", 50)
	viper.SetDefault("chat_limit", 50)

	// Config file is optional
	_ = viper.ReadInConfig()

	cfg := &Config{
		ServerURL:       viper.GetString("server_url"),
		Password:        viper.GetString("password"),
		PollIntervalSec: viper.GetInt("poll_interval_sec"),
		MessageLimit:    viper.GetInt("message_limit"),
		ChatLimit:       viper.GetInt("chat_limit"),
	}

	if cfg.ServerURL == "" || cfg.Password == "" {
		return nil, fmt.Errorf("BB_SERVER_URL and BB_PASSWORD environment variables are required")
	}

	return cfg, nil
}
