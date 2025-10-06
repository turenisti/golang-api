package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ServiceCode string
	AppPort     string
	Environment string

	// Database
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBType string

	// Logging
	LogLevel  string
	LogFormat string

	// Scheduler
	SchedulerEnabled bool
}

var Config AppConfig

func LoadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn().Err(err).Msg("No config file found, using environment variables")
	}

	Config = AppConfig{
		ServiceCode:      viper.GetString("SERVICE_CODE"),
		AppPort:          viper.GetString("APP_PORT"),
		Environment:      viper.GetString("ENVIRONMENT"),
		DBHost:           viper.GetString("DB_HOST"),
		DBPort:           viper.GetString("DB_PORT"),
		DBUser:           viper.GetString("DB_USER"),
		DBPass:           viper.GetString("DB_PASS"),
		DBName:           viper.GetString("DB_NAME"),
		DBType:           viper.GetString("DB_TYPE"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		LogFormat:        viper.GetString("LOG_FORMAT"),
		SchedulerEnabled: viper.GetBool("SCHEDULER_ENABLED"),
	}

	log.Info().Msg("Configuration loaded successfully")
}
