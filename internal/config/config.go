package config

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Secret        string
	DbUrl         string
	EmailPassword string
}

var CFG *Config

func LoadConfig() {
	err := godotenv.Load(filepath.Join(".", ".env"))

	if err != nil {
		log.Println("Error loading .env file")
	}

	if CFG == nil {
		CFG = &Config{
			Port:          os.Getenv("PORT"),
			Secret:        os.Getenv("SECRET"),
			DbUrl:         os.Getenv("DB_URL"),
			EmailPassword: os.Getenv("EMAIL_PASSWORD"),
		}
	}

}
