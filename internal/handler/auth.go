package handler

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yuusufyan/go-common/pkg/utils"
	"gorm.io/gorm"

	"backend-sevima/internal/model"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user against the database and returns a JWT
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user model.User
	// Find user by email and preload tenant to ensure it exists
	if err := h.db.Preload("Tenant").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !user.IsActive || !user.Tenant.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User or Tenant is inactive",
		})
	}

	// Create JWT token
	secret := os.Getenv("JWT_TOKEN_SECRET")
	if secret == "" {
		secret = "TunaKaleng!!!!"
	}

	claims := jwt.MapClaims{
		"sub":       user.ID,
		"tenant_id": user.TenantID,
		"email":     user.Email,
		"role":      user.Role,
		"type":      "access",
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Login successful",
		"token":     tokenString,
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
	})
}
