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

	// Set all environment variables from program context
	// If the same variable will be in .env file it will NOT be overwritten
	c.BitlyOAuthToken = viper.GetString("BITLY_OAUTH_TOKEN")

	// Reading public config file
	if err := ReadAndParseFromFile(path, &c); err != nil {
		panic(err)
	}

	// Read .env file in root directory
	// If file not exists, then display warning and continue program
	envPath := "./.env"
	if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf(".env file was not found in %s\n", envPath)
	} else {
		if err = ReadAndParseFromFile(envPath, &c); err != nil {
			panic(err)
		}
	}

	// Check required variables manually
	if c.BitlyOAuthToken == "" {
		panic("BITLY_OAUTH_TOKEN is not set")
	}

	return &c
}
