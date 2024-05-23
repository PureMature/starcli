package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration with the given path.
func InitConfig(configPath string) error {
	if configPath != "" {
		// use the provided config file
		viper.SetConfigFile(configPath)
	} else {
		// try to find home and config directories
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config", AppName)

		// search config in current and config directory with name app name (without extension)
		viper.AddConfigPath(".")
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// set default values
	viper.SetTypeByDefaultValue(true)
	// TODO: set default values

	// read in environment variables that match
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(AppName)

	// read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			return err
		}
	}

	// done
	return nil
}
