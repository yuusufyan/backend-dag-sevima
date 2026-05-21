package routing

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend-sevima/internal/handlers"
	"backend-sevima/internal/services"
)

// SetupAuthRoutes sets up the authentication routes
func SetupAuthRoutes(api fiber.Router, db *gorm.DB) {
	// Initialize services
	authService := services.NewAuthService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Routes
	authGroup := api.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/register", authHandler.Register)
}
