package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yuusufyan/go-common/response"

	"backend-sevima/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user authentication request
// @Summary User Login
// @Description Authenticates a user against the database and returns a JWT
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login Request"
// @Success 200 {object} response.Response[services.LoginResponse]
// @Failure 400 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	result, err := h.authService.Login(req)
	if err != nil {
		return response.RespondWithError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Login successful", result)
}

// Register handles tenant and user registration
// @Summary User Registration
// @Description Creates a new tenant and a new admin user within a transaction
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Register Request"
// @Success 201 {object} response.Response[services.RegisterResponse]
// @Failure 400 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	result, err := h.authService.Register(req)
	if err != nil {
		return response.RespondWithError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "Tenant and User created successfully", result)
}
