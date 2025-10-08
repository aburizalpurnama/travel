package config

import "github.com/joho/godotenv"

// LoadConfig memuat environment variable dari file .env
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
