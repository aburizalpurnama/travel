package driver

import (
	"fmt"

	"github.com/aburizalpurnama/travel/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgreDatabase return gorm dbmap object with postgre options param
func NewPostgreDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getLoggerLevel(cfg.DBLogLevel)),
	})
	if err != nil {
		return nil, err
	}

	pgsqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = pgsqlDB.Ping()
	if err != nil {
		return nil, err
	}

	pgsqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	pgsqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	pgsqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)

	return db, nil
}

func getLoggerLevel(v string) logger.LogLevel {
	switch v {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
