package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("API_BASE_URL")
	if dbURL == "" {
		log.Fatal("API_BASE_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform != "dev" && platform != "prod" {
		log.Fatal("PLATFORM must be set to either dev or prod")
	}

	startRepl()
}
