package config

import (
	"fmt"
	"github.com/spf13/viper"
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
