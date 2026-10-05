package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"naughtfound.github.io/go2phd/handlers"
	"naughtfound.github.io/go2phd/repositories"
	"naughtfound.github.io/go2phd/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

	if spreadsheetID == "" {
		log.Fatalf("SPREADSHEET_ID environment variable is missing")
	}

	uniRepo := repositories.NewUniversityRepository()

	var sheetsRepo *repositories.SheetsRepository
	var err error

	if refreshToken != "" && clientID != "" && clientSecret != "" {
		log.Println("Initializing Google Sheets Repository using Refresh Token from environment")
		sheetsRepo, err = repositories.NewSheetsRepositoryFromRefreshToken(clientID, clientSecret, refreshToken, spreadsheetID)
	} else {
		log.Println("Initializing Google Sheets Repository using credentials.json")
		sheetsRepo, err = repositories.NewSheetsRepository("credentials.json", spreadsheetID)
	}

	if err != nil {
		log.Fatalf("Failed to initialize Google Sheets repository: %v", err)
	}

	uniService := services.NewUniversityService(uniRepo)
	posService := services.NewPositionService(sheetsRepo, uniRepo)
	authService := services.NewAuthService(clientID, clientSecret, redirectURL)

	uniHandler := handlers.NewUniversityHandler(uniService)
	posHandler := handlers.NewPositionHandler(posService)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	authGroup := r.Group("/auth")
	authHandler.SetupGroup(authGroup)

	universityGroup := r.Group("/universities")
	uniHandler.SetupGroup(universityGroup)

	positionGroup := r.Group("/positions")
	posHandler.SetupGroup(positionGroup)

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
