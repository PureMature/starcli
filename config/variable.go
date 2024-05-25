package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	emptyStr string
)

// SetDefaults sets the default values in Viper for the configuration.
func SetDefaults() {
	host, _ := os.Hostname()
	viper.SetDefault("host_name", host)
	viper.SetDefault("resend_api_key", emptyStr)
	viper.SetDefault("sender_domain", emptyStr)
}

// GetHostname returns the host name from the configuration.
func GetHostname() string {
	return viper.GetString("host_name")
}

// GetResendAPIKey returns the resend API key from the configuration.
func GetResendAPIKey() string {
	return viper.GetString("resend_api_key")
}

// GetSenderDomain returns the email sender domain from the configuration.
func GetSenderDomain() string {
	return viper.GetString("sender_domain")
}
