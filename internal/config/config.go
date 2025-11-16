package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string `env:"APP_ENV"       envDefault:"development"`
	AppName     string `env:"APP_NAME"      envDefault:"travel-api"` // TODO: default value ganti dengan nama app pada builder
	AppLogLevel string `env:"APP_LOG_LEVEL"`
	ServerPort  int    `env:"SERVER_PORT"   envDefault:"3000"`

	DBHost            string        `env:"DB_HOST,required"`
	DBPort            int           `env:"DB_PORT,required"`
	DBUsername        string        `env:"DB_USERNAME,required"`
	DBPassword        string        `env:"DB_PASSWORD,required"`
	DBName            string        `env:"DB_NAME,required"`
	DBSSLMode         string        `env:"DB_SSLMODE"           envDefault:"disable"`
	DBTimezone        string        `env:"DB_TIMEZONE"          envDefault:"Asia/Jakarta"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS"    envDefault:"10"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS"    envDefault:"100"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	DBLogLevel        string        `env:"DB_LOG_LEVEL"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"       envDefault:"0"`

	JwtSecret            string `env:"JWT_SECRET"`
	JwtExpirationMinutes int    `env:"JWT_EXPIRATION_MINUTES" envDefault:"60"`

	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:5173,http://localhost:3000"`

	MailgunApiKey   string `env:"MAILGUN_API_KEY"`
	MailgunDomain   string `env:"MAILGUN_DOMAIN"`
	MailSenderEmail string `env:"MAILGUN_SENDER_EMAIL"`

	Tracing struct {
		Enabled  bool   `env:"TRACING_ENABLED" envDefault:"false"`
		Exporter string `env:"TRACING_EXPORTER" envDefault:"stdout"` // "stdout" or "otlp"

		// OTLP General Config
		OtlpEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"localhost:4317"`
		OtlpInsecure bool   `env:"OTEL_EXPORTER_OTLP_INSECURE" envDefault:"true"`

		// (Opsional) Vendor-specific headers
		OtlpHeaders string `env:"OTEL_EXPORTER_OTLP_HEADERS"` // e.g., "Authentication=Bearer <key>,Datadog-Meta-Tracer-Version=v1"
	}
}

func LoadConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "development" || appEnv == "" {
		err := godotenv.Load()
		if err != nil && !os.IsNotExist(err) {
			// Beri peringatan jika ada error lain (misal: format .env salah),
			// tapi abaikan jika error-nya hanya "file tidak ditemukan".
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
