package routing

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"


	_ "backend-sevima/docs" // Add this for swagger initialization
	"github.com/yuusufyan/go-common/pkg/middleware/swagger"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Setup Swagger from go-common
	swagger.SetupSwagger(app)

	api := app.Group("/api/v1")

	// Setup modular routes
	SetupAuthRoutes(api, db)
	SetupTenantRoutes(api, db)
}
