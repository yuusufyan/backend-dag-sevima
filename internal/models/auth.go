package model

import (
	"github.com/yuusufyan/go-common/pkg/database"
)

// Tenant represents a multi-tenant client or workspace.
type Tenant struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	database.AuditModel
}

// User represents a user that can log in and belong to a specific tenant.
type User struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID string `gorm:"type:uuid;index;not null"`
	Email    string `gorm:"size:255;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:50;default:'member'"`
	IsActive bool   `gorm:"default:true"`
	database.AuditModel

	Tenant Tenant `gorm:"foreignKey:TenantID"`
}
