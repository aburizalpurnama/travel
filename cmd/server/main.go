package main

import (
	"fmt"
	"log"

	"github.com/aburizalpurnama/travel/internal/app/database"
	"github.com/aburizalpurnama/travel/internal/app/domain/user"
	"github.com/aburizalpurnama/travel/internal/app/router"
	"github.com/aburizalpurnama/travel/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	routerOpts := injectDependencies(db)

	app := fiber.New()
	app.Use(logger.New())

	router.SetupRoutesV1(app, routerOpts)

	port := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Fatal(app.Listen(port))
}

func injectDependencies(db *gorm.DB) *router.Option {
	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	return &router.Option{
		UserHandler: userHandler,
	}
}
