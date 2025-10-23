package database

import (
	"database/sql"
	"fmt"

	"github.com/aburizalpurnama/travel/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimezone,
	)

	logLevel := getLogLevel(cfg.AppEnv)
	if cfg.DBLogLevel != "" {
		logLevel = convertLogLevel(cfg.DBLogLevel)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// GORM AutoMigrate hanya untuk development cepat.
	// Di produksi, sebaiknya gunakan tool migrasi terpisah (seperti GORM Migrate atau Goose).

	fmt.Println("Database connection successful")
	return db, nil
}

func NewNativeSQL(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimezone,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connection successful")
	return db, nil
}

func getLogLevel(v string) logger.LogLevel {
	switch v {
	case "development":
		return logger.Info
	case "staging":
		return logger.Info
	case "production":
		return logger.Warn
	default:
		return logger.Info
	}
}

func convertLogLevel(v string) logger.LogLevel {
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
