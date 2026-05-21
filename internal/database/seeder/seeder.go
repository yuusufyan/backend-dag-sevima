package seeder

import (
	"errors"

	model "backend-sevima/internal/models"
	"github.com/yuusufyan/go-common/pkg/logger"
	"github.com/yuusufyan/go-common/pkg/utils"
	"gorm.io/gorm"
)

// SeedSuperAdmin ensures that a default tenant and an Admin user exist in the database.
func SeedSuperAdmin(db *gorm.DB, appLogger logger.Logger) {
	// Check if master tenant exists
	var tenant model.Tenant
	err := db.Where("name = ?", "System").First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create System Tenant
			tenant = model.Tenant{
				Name:     "System",
				IsActive: true,
			}
			if err := db.Create(&tenant).Error; err != nil {
				appLogger.WithError(err).Error("Failed to seed System Tenant")
				return
			}
			appLogger.Info("Seeded System Tenant")
		} else {
			appLogger.WithError(err).Error("Error checking System Tenant")
			return
		}
	}

	// Check if admin user exists
	var admin model.User
	err = db.Where("email = ?", "admin@sevima.com").First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, errHash := utils.HashPassword("admin123")
			if errHash != nil {
				appLogger.WithError(errHash).Error("Failed to hash admin password")
				return
			}
			
			// Create Admin User
			admin = model.User{
				TenantID: tenant.ID,
				Email:    "admin@sevima.com",
				Password: hashedPassword,
				Role:     "Admin",
				IsActive: true,
			}
			if err := db.Create(&admin).Error; err != nil {
				appLogger.WithError(err).Error("Failed to seed Admin User")
				return
			}
			appLogger.Info("Seeded Admin User (admin@sevima.com / admin123)")
		} else {
			appLogger.WithError(err).Error("Error checking Admin User")
			return
		}
	}
}
