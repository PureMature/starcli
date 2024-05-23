package config

import (
	"os"

	"github.com/spf13/viper"
)

// SetDefaults sets the default values in Viper for the configuration.
func SetDefaults() {
	host, _ := os.Hostname()
	viper.SetDefault("host_name", host)
}

// GetHostname returns the host name from the configuration.
func GetHostname() string {
	return viper.GetString("host_name")
}
