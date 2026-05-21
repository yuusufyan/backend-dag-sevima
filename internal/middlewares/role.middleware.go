package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole checks if the authenticated user has one of the required roles.
// It assumes AuthGuard has already run and placed "role" in c.Locals.
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: No role assigned",
			})
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Insufficient permissions",
		})
	}
}
