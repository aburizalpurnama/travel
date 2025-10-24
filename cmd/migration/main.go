package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"

	// Impor driver postgres
	_ "github.com/jackc/pgx/v5/stdlib"

	// Impor file .go migrasi Anda
	"github.com/aburizalpurnama/travel/internal/app/database"
	_ "github.com/aburizalpurnama/travel/internal/app/database/migration"
	"github.com/aburizalpurnama/travel/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.NewNativeSQL(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %s", err.Error())
		}
	}()

	// Ambil argumen (up, down, status)
	if len(os.Args) < 2 {
		log.Fatalf("Missing command: 'up', 'down', or 'status'")
	}

	command := os.Args[1]
	migrationDir := "internal/app/database/migration"

	// Jalankan perintah goose
	if err := goose.RunContext(context.Background(), command, db, migrationDir); err != nil {
		log.Fatalf("Goose run failed for command '%s': %v", command, err)
	}

	fmt.Printf("Goose command '%s' ran successfully.\n", command)
}
