package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	instance *Env
	once     sync.Once
)

type Env struct {
	ENV               string
	SERVER_PORT       string
	DATABASE_DRIVER   string
	DATABASE_HOST     string
	DATABASE_USER     string
	DATABASE_NAME     string
	DATABASE_PASSWORD string
	DATABASE_PORT     string
	DATABASE_SCHEMA   string

	CLOUD_ENV          string
	CLOUD_KEY          string
	CLOUD_SECRET       string
	CLOUD_REGION       string
	CLOUD_DISABLED_SSL bool
	CLOUD_HOST         string
	CLOUD_BUCKET       string
	CLOUD_HOST_BUCKET  string
}

func Load() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}

		instance = &Env{
			ENV:               getEnv("ENV", "development"),
			SERVER_PORT:       getEnv("SERVER_PORT", ""),
			DATABASE_HOST:     getEnv("DATABASE_HOST", ""),
			DATABASE_USER:     getEnv("DATABASE_USER", ""),
			DATABASE_NAME:     getEnv("DATABASE_NAME", ""),
			DATABASE_PASSWORD: getEnv("DATABASE_PASSWORD", ""),
			DATABASE_PORT:     getEnv("DATABASE_PORT", ""),
			DATABASE_SCHEMA:   getEnv("DATABASE_SCHEMA", ""),

			CLOUD_ENV:          getEnv("CLOUD_ENV", "aws"),
			CLOUD_KEY:          getEnv("CLOUD_KEY", "AKIAUALMXV557XUQV26Y"),
			CLOUD_SECRET:       getEnv("CLOUD_SECRET", ""),
			CLOUD_REGION:       getEnv("CLOUD_REGION", "sa-east-1"),
			CLOUD_DISABLED_SSL: getEnvAsBool("CLOUD_DISABLED_SSL", true),
			CLOUD_HOST:         getEnv("CLOUD_HOST", "https://s3.sa-east-1.amazonaws.com"),
			CLOUD_BUCKET:       getEnv("CLOUD_BUCKET", "market-prd"),
			CLOUD_HOST_BUCKET:  getEnv("CLOUD_HOST_BUCKET", "https://market-prd.s3.sa-east-1.amazonaws.com"),
		}
	})

}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" {
			return true
		}
		if value == "false" || value == "0" {
			return false
		}
	}
	return fallback
}

func Get() *Env {
	if instance == nil {
		panic("Envuration not loaded. Call config.Load() first.")
	}
	return instance
}
