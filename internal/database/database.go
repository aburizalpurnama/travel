package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// GORM AutoMigrate hanya untuk development cepat.
	// Di produksi, sebaiknya gunakan tool migrasi terpisah (seperti GORM Migrate atau Goose).
	// Karena Anda sudah memiliki DDL, jalankan DDL itu langsung di database.
	// db.AutoMigrate(&domain.User{})

	fmt.Println("Database connection successful")
	return db, nil
}
