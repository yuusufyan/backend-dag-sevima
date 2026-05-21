package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yuusufyan/go-common/pkg/database"
	"gorm.io/gorm"
)

// Tenant represents a multi-tenant client or workspace.
type Tenant struct {
	ID          uuid.UUID `gorm:"type:uuid;column:id" json:"id"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"default:true"`
	database.AuditModel
}

// User represents a user that can log in and belong to a specific tenant.
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;column:id" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;index;not null"`
	Email     string         `gorm:"size:255;unique;not null"`
	Password  string         `gorm:"size:255;not null"`
	Role      string         `gorm:"size:50;default:'member'"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Tenant Tenant `gorm:"foreignKey:TenantID;null"`
}
