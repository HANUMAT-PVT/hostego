package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables.")
	}
}

// GetEnv fetches an environment variable with a default fallback
func GetEnv(key string) string {
	LoadEnv()
	appEnv, exists := os.LookupEnv("APP_ENV")
	fmt.Println(appEnv, "ap env")

	value, exists := os.LookupEnv(key + appEnv)
	if !exists {
		return ""
	}
	return value
}
