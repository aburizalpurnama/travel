package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds all the application configuration values loaded from environment variables.
type Config struct {
	// Application Settings
	AppEnv      string `env:"APP_ENV"       envDefault:"development"`
	AppName     string `env:"APP_NAME"      envDefault:"travel-api"` // TODO: Replace default value with the actual app name during build
	AppLogLevel string `env:"APP_LOG_LEVEL"`
	ServerPort  int    `env:"SERVER_PORT"   envDefault:"3000"`

	// Database Configuration
	DatabaseURL       string        `env:"DATABASE_URL"`
	DBHost            string        `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort            int           `env:"DB_PORT" envDefault:"5432"`
	DBUsername        string        `env:"DB_USERNAME"`
	DBPassword        string        `env:"DB_PASSWORD"`
	DBName            string        `env:"DB_NAME"`
	DBSSLMode         string        `env:"DB_SSLMODE"           envDefault:"disable"`
	DBTimezone        string        `env:"DB_TIMEZONE"          envDefault:"Asia/Jakarta"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS"    envDefault:"10"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS"    envDefault:"100"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	DBLogLevel        string        `env:"DB_LOG_LEVEL"`

	// Redis Configuration
	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"       envDefault:"0"`

	// Security Configuration
	JwtSecret            string `env:"JWT_SECRET"`
	JwtExpirationMinutes int    `env:"JWT_EXPIRATION_MINUTES" envDefault:"60"`
	CORSAllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS"   envDefault:"http://localhost:5173,http://localhost:3000"`

	// Email Service Configuration (Mailgun)
	MailgunApiKey   string `env:"MAILGUN_API_KEY"`
	MailgunDomain   string `env:"MAILGUN_DOMAIN"`
	MailSenderEmail string `env:"MAILGUN_SENDER_EMAIL"`

	// OpenTelemetry Tracing Configuration
	Tracing struct {
		Enabled  bool   `env:"TRACING_ENABLED"  envDefault:"false"`
		Exporter string `env:"TRACING_EXPORTER" envDefault:"stdout"` // Options: "stdout", "otlp"

		// OTLP General Config
		OtlpEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
		OtlpInsecure bool   `env:"OTEL_EXPORTER_OTLP_INSECURE" envDefault:"true"`

		// Optional Vendor-specific headers (e.g., "Authentication=Bearer <key>,...")
		OtlpHeaders string `env:"OTEL_EXPORTER_OTLP_HEADERS"`
	}
}

// LoadConfig reads configuration from .env files and environment variables.
// It prioritizes existing environment variables over .env file values.
func LoadConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	// Load .env file only in development or if APP_ENV is not set
	if appEnv == "development" || appEnv == "" {
		err := godotenv.Load()
		if err != nil && !os.IsNotExist(err) {
			// Warn on errors other than "file not found" (e.g., invalid .env format)
			log.Printf("Warning: Non-critical error loading .env file: %v", err)
		}
	}

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
