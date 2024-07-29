package config

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	Secret string
	DbUrl  string
}

// Singleton Pattern
var lock = &sync.Mutex{}

// Will change to package-level variable later
var cfg *Config

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if cfg == nil {
		lock.Lock()
		defer lock.Unlock()

		if cfg == nil {
			cfg = &Config{
				Port:   os.Getenv("PORT"),
				Secret: os.Getenv("SECRET"),
				DbUrl:  os.Getenv("DB_URL"),
			}
		}
	}
	return cfg
}
