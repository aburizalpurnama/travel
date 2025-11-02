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
	// Tentukan jumlah percobaan dan waktu tunggu
	maxRetries    = 10              // Coba 10 kali
	retryInterval = 5 * time.Second // Tunggu 5 detik antar percobaan
)

// retryOperation adalah helper generik untuk menjalankan operasi dengan retry
func retryOperation(operation func() error) error {
	var lastError error

	for i := 1; i <= maxRetries; i++ {
		lastError = operation() // Jalankan operasi (misal: connect atau ping)
		if lastError == nil {
			return nil // Sukses, tidak perlu error
		}

		// Jika gagal, catat log dan tunggu sebelum mencoba lagi
		log.Printf("[DB Connect] Percobaan %d/%d gagal: %v", i, maxRetries, lastError)
		if i < maxRetries {
			log.Printf("Mencoba lagi dalam %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	// Jika loop selesai, berarti semua percobaan gagal
	return fmt.Errorf("gagal terhubung ke database setelah %d percobaan: %w", maxRetries, lastError)
}

// --- FUNGSI ANDA YANG SUDAH DIPERBAIKI ---

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

	// Gunakan helper retry untuk membungkus gorm.Open
	// gorm.Open sudah otomatis melakukan "ping" saat koneksi
	op := func() error {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		return err
	}

	err = retryOperation(op)
	if err != nil {
		return nil, err // Kembalikan error terakhir jika semua percobaan gagal
	}

	log.Println("✅ Koneksi database GORM berhasil dibuat!")
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

	// Ganti "pgx" dengan "pgx/v5/stdlib" jika Anda menggunakan pgx v5
	// sql.Open TIDAK membuka koneksi, jadi ini tidak akan gagal.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal mem-parsing DSN: %w", err)
	}

	// Gunakan helper retry untuk membungkus db.Ping()
	// db.Ping() adalah yang BENAR-BENAR menguji koneksi.
	op := func() error {
		return db.Ping()
	}

	err = retryOperation(op)
	if err != nil {
		err = db.Close()      // Tutup koneksi yang gagal
		return nil, err // Kembalikan error terakhir jika semua percobaan gagal
	}

	log.Println("✅ Koneksi database (Native SQL) berhasil dibuat!")
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
