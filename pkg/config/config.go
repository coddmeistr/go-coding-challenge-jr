package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type BiltyOAuth struct {
	Token string
}

type ServerConfig struct {
	Port       int `mapstructure:"port"`
	BiltyOAuth BiltyOAuth
}

// MustLoadByPath load envs and marshaling config file in given path
//
// It panics on any error
func MustLoadByPath(path string) *ServerConfig {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var c ServerConfig

	err = viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	_ = viper.GetString("BITLY_OAUTH_LOGIN")
	c.BiltyOAuth.Token = viper.GetString("BITLY_OAUTH_TOKEN")

	return &c
}
