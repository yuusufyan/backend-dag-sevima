package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yuusufyan/go-common/pkg/utils"
)

// AuthGuard protects endpoints by ensuring a valid JWT is present in the Authorization header.
// It parses the token, verifies it, and extracts the tenant_id and user_id into the context.
func AuthGuard(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing Authorization header",
		})
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid Authorization header format",
		})
	}

	tokenStr := parts[1]
	secret := os.Getenv("JWT_TOKEN_SECRET")
	if secret == "" {
		secret = "TunaKaleng!!!!" // Fallback matching .env for development
	}

	// Token type is usually "access" or similar, we'll assume "access" or omit if the token doesn't have it.
	// We'll pass "access" as expected tokenType, assuming go-common validates it against claims["type"]
	claims, err := utils.VerifyToken(tokenStr, secret, "access")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}

	// Extract Tenant ID from claims
	tenantID, ok := claims["tenant_id"].(string)
	if !ok || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing tenant_id in token claims",
		})
	}

	// Store the extracted values in locals for downstream handlers
	c.Locals("tenant_id", tenantID)

	// Additionally extract user_id (sub)
	if sub, ok := claims["sub"].(string); ok {
		c.Locals("user_id", sub)
	}

	// Extract role
	if role, ok := claims["role"].(string); ok {
		c.Locals("role", role)
	}

	return c.Next()
}
