package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

// ReadAndParseFromFile takes config file and loads this config file in viper
//
// If dest parameter is not nil, then function will try to unmarshal config file into dest
func ReadAndParseFromFile(configFile string, dest any) error {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}
	if dest != nil {
		if err := viper.Unmarshal(dest); err != nil {
			return fmt.Errorf("unable to decode into struct, %v", err)
		}
	}

	return nil
}

// LoadEnvs helper functions that loads envs in viper with viper.AutomaticEnv
// and loads envs from given file from given path.
//
// If file not exists or some error occurred, then it prints warning message.
func LoadEnvs(path string) {
	viper.AutomaticEnv()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("file with envs was not found in %s\n", path)
		return
	}
	err := ReadAndParseFromFile(path, nil)
	if err != nil {
		log.Printf("couldn't load envs from file in: %s, err: %v", path, err)
	}
}
