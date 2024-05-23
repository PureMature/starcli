package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func InitConfig(configPath string) {
	viper.SetConfigName("config") // Default configuration file name (without extension)
	viper.SetConfigType("yaml")   // Default configuration file format
	viper.AddConfigPath(".")      // Look for the config in the current directory
	viper.AddConfigPath("$HOME")  // Look for the config in the user's home directory

	// If a specific config file is provided, use it
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found, using defaults and environment variables")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	// Allow environment variables to override configuration settings
	viper.AutomaticEnv()
}
