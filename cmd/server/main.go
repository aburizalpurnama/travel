package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aburizalpurnama/travel/internal/config"
	"github.com/aburizalpurnama/travel/internal/database"
	"github.com/aburizalpurnama/travel/internal/router"
	"github.com/aburizalpurnama/travel/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Inisialisasi database
	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Inisialisasi arsitektur (Dependency Injection)
	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	// Inisialisasi Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Setup rute
	router.SetupRoutes(app, userHandler)

	// Jalankan server
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Fatal(app.Listen(port))
}
