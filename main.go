package main

import (
	"log"
	"os"
	"time"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
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

	session := digitalshelfapi.Session{
		DSAPIClient: digitalshelfapi.NewClient(time.Second * 10),
		Base_url:    dbURL,
		Platform:    platform,
	}

	startRepl(&session)
}
