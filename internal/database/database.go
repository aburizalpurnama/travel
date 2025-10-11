package database

import (
	"fmt"

	"github.com/aburizalpurnama/travel/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// GORM AutoMigrate hanya untuk development cepat.
	// Di produksi, sebaiknya gunakan tool migrasi terpisah (seperti GORM Migrate atau Goose).

	fmt.Println("Database connection successful")
	return db, nil
}
