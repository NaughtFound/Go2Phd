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

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	uniRepo := repositories.NewUniversityRepository()

	uniService := services.NewUniversityService(uniRepo)
	posService := services.NewPositionService(uniRepo, clientID, clientSecret)
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
