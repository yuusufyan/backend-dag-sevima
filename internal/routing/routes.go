package routing

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend-sevima/internal/handlers"
	middleware "backend-sevima/internal/middlewares"
	"backend-sevima/internal/services"

	_ "backend-sevima/docs" // Add this for swagger initialization
	"github.com/yuusufyan/go-common/pkg/middleware/swagger"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Setup Swagger from go-common
	swagger.SetupSwagger(app)

	api := app.Group("/api/v1")

	// Initialize services
	authService := services.NewAuthService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Public Routes
	authGroup := api.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/register", authHandler.Register)

	// Protected Routes Group
	// Apply Custom Auth Guard (verifies JWT & extracts Tenant ID)
	protected := api.Group("/")
	protected.Use(middleware.AuthGuard)

	// Example protected route:
	// protected.Get("/profile", handlers.GetProfile)
}
