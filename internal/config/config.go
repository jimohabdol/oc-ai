package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAIKey      string `mapstructure:"openai_key"`
	DefaultModel   string `mapstructure:"default_model"`
	ConfirmExecute bool   `mapstructure:"confirm_execute"`
	HistoryLimit   int    `mapstructure:"history_limit"`
}

func LoadConfig() (*Config, error) {
	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Search paths
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}
	viper.AddConfigPath(filepath.Join(configDir, "oc-ai"))
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("OC_AI")
	viper.AutomaticEnv()
	viper.BindEnv("openai_key")
	viper.BindEnv("default_model")

	// Defaults
	viper.SetDefault("default_model", "gpt-4-turbo")
	viper.SetDefault("confirm_execute", true)
	viper.SetDefault("history_limit", 100)

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file found but another error occurred
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Fallback to environment variable if not set in config
	if cfg.OpenAIKey == "" {
		cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	}

	return &cfg, nil
}
