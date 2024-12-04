package config

import (
	"github.com/spf13/viper"
	"log"
)

// LoadConfig loads the configuration file
func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

// GetString retrieves a string value from the configuration
func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}
