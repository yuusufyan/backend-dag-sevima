package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthGuard ensures that the request has a valid Tenant ID.
// In a real scenario, this would parse a JWT token and extract the Tenant ID from the claims.
// For now, we will extract it from a custom header "X-Tenant-ID" to simulate auth.
func AuthGuard(c *fiber.Ctx) error {
	tenantID := c.Get("X-Tenant-ID")
	
	if tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing X-Tenant-ID header",
		})
	}

	// Store the tenant ID in locals for downstream handlers to use
	c.Locals("tenant_id", tenantID)

	// Additional User context can be stored here when JWT is implemented
	// c.Locals("user_id", "user-uuid")

	return c.Next()
}
