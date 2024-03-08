package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type ServerConfig struct {
	Port            int    `mapstructure:"port"`
	BitlyOAuthToken string `mapstructure:"BITLY_OAUTH_TOKEN"`
}

// MustLoadByPath load envs and marshaling config file in given path
//
// It panics on any error
func MustLoadByPath(path string) *ServerConfig {
	viper.AutomaticEnv()

	var c ServerConfig

	// Reading public config file
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	// Read .env file in root directory
	// If file not exists, skip viper's parsing logic
	envPath := "./.env"
	if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("envs file not found with this filename: %s\n", envPath)

		c.BitlyOAuthToken = viper.GetString("BITLY_OAUTH_TOKEN")
		goto skipEnvFile
	}
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

skipEnvFile:

	if c.BitlyOAuthToken == "" {
		panic("BITLY_OAUTH_TOKEN is not set")
	}

	return &c
}
