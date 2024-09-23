package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func GetEnv(key string, fallback string) string {
	if value, exists := viper.Get(key).(string); exists {
		return value
	}

	return fallback
}