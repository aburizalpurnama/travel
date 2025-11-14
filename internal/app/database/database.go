package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aburizalpurnama/travel/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetries    = 10              // Number of connection attempts
	retryInterval = 5 * time.Second // Wait time between attempts
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

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var db *gorm.DB
	var err error

	// Define the connection operation
	op := func() error {
		// gorm.Open attempts to connect and ping the DB automatically
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		return err
	}

	err = retryOperation(op)
	if err != nil {
		return nil, err
	}

	log.Println("✅ GORM database connection established!")
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

	// sql.Open only prepares the DSN; it does not connect.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Define the ping operation to validate the connection
	op := func() error {
		return db.Ping()
	}

	err = retryOperation(op)
	if err != nil {
		db.Close() // Clean up the failed connection
		return nil, err
	}

	log.Println("✅ Native SQL database connection established!")
	return db, nil
}

// retryOperation provides a generic wrapper to execute an operation with retry logic.
func retryOperation(operation func() error) error {
	var lastError error

	for i := 1; i <= maxRetries; i++ {
		lastError = operation()
		if lastError == nil {
			return nil // Success
		}

		log.Printf("[DB Connect] Attempt %d/%d failed: %v", i, maxRetries, lastError)
		if i < maxRetries {
			log.Printf("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, lastError)
}

// getLogLevel sets GORM log level based on application environment.
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

// convertLogLevel sets GORM log level based on a specific config string.
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
