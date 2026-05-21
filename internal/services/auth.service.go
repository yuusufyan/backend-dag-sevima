package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yuusufyan/go-common/pkg/apperror"
	"github.com/yuusufyan/go-common/pkg/utils"
	"gorm.io/gorm"

	model "backend-sevima/internal/models"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	var user model.User
	if err := s.db.Preload("Tenant").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.Unauthorized("Invalid credentials")
		}
		return nil, apperror.InternalServer("Database error")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, apperror.Unauthorized("Invalid credentials")
	}

	if !user.IsActive || !user.Tenant.IsActive {
		return nil, apperror.Unauthorized("User or Tenant is inactive")
	}

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
		return nil, apperror.InternalServer("Failed to generate token")
	}

	return &LoginResponse{
		Token:    tokenString,
		UserID:   user.ID.String(),
		TenantID: user.TenantID.String(),
	}, nil
}

type RegisterRequest struct {
	TenantName string `json:"tenant_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type RegisterResponse struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, apperror.BadRequest("Tenant name, email, and password are required")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.InternalServer("Failed to hash password")
	}

	tx := s.db.Begin()

	tenant := model.Tenant{
		Name:     req.TenantName,
		IsActive: true,
	}
	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return nil, apperror.InternalServer("Failed to create tenant")
	}

	user := model.User{
		TenantID: tenant.ID,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "admin",
		IsActive: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, apperror.InternalServer("Failed to create user. Email might already exist.")
	}

	tx.Commit()

	return &RegisterResponse{
		TenantID: tenant.ID.String(),
		UserID:   user.ID.String(),
	}, nil
}
