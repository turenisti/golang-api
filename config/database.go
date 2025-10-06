package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// Build MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.DBUser,
		Config.DBPass,
		Config.DBHost,
		Config.DBPort,
		Config.DBName,
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Use custom logger
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Connection pooling
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().
		Str("type", Config.DBType).
		Str("host", Config.DBHost).
		Str("database", Config.DBName).
		Msg("Database connected successfully")
}
