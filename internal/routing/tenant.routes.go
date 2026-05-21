package routing

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"backend-sevima/internal/handlers"
	"backend-sevima/internal/services"
)

// SetupTenantRoutes sets up the tenant CRUD routes
func SetupTenantRoutes(api fiber.Router, db *gorm.DB) {
	// Initialize services
	tenantService := services.NewTenantService(db)

	// Initialize handlers
	tenantHandler := handlers.NewTenantHandler(tenantService)

	// Protected Routes Group
	tenantGroup := api.Group("/tenants")
	// tenantGroup.Use(middleware.AuthGuard) // Require JWT authentication

	// Routes
	tenantGroup.Get("/", tenantHandler.GetAll)
	tenantGroup.Get("/:id", tenantHandler.GetByID)

	// Protected routes requiring Admin/SuperAdmin privileges
	// adminGuard := middleware.RequireRole("admin", "superadmin", "Admin", "SuperAdmin")
	tenantGroup.Post("/", tenantHandler.Create)
	tenantGroup.Put("/:id", tenantHandler.Update)
	tenantGroup.Delete("/:id", tenantHandler.Delete)
}
