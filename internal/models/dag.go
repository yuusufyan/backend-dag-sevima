package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yuusufyan/go-common/pkg/database"
	"gorm.io/datatypes"
)

// DagTemplate stores the definition of a DAG.
type DagTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;column:create_by" json:"create_by"`
	TenantID    string    `gorm:"index;not null"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	// Definition contains the raw JSON structure of the DAG (nodes and edges)
	Definition datatypes.JSON `gorm:"type:jsonb"`
	database.AuditModel
}

// DagExecution tracks a single execution instance of a DagTemplate.
type DagExecution struct {
	ID            uuid.UUID      `gorm:"type:uuid;column:create_by" json:"create_by"`
	TenantID      string         `gorm:"index;not null"`
	DagTemplateID string         `gorm:"type:uuid;index;not null"`
	Status        string         `gorm:"size:50;not null"` // PENDING, RUNNING, SUCCESS, FAILED
	Input         datatypes.JSON `gorm:"type:jsonb"`
	Output        datatypes.JSON `gorm:"type:jsonb"`
	Error         string         `gorm:"type:text"`
	StartedAt     *time.Time
	CompletedAt   *time.Time
	database.AuditModel

	Template DagTemplate `gorm:"foreignKey:DagTemplateID"`
}

// TaskInstance tracks the execution of a single Task within a DagExecution.
type TaskInstance struct {
	ID             uuid.UUID      `gorm:"type:uuid;column:create_by" json:"create_by"`
	TenantID       string         `gorm:"index;not null"`
	DagExecutionID string         `gorm:"type:uuid;index;not null"`
	NodeID         string         `gorm:"size:255;not null"` // ID of the node in the DAG definition
	TaskType       string         `gorm:"size:100;not null"` // HTTP, SCRIPT, DELAY, etc.
	Status         string         `gorm:"size:50;not null"`  // PENDING, RUNNING, SUCCESS, FAILED, SKIPPED
	Input          datatypes.JSON `gorm:"type:jsonb"`
	Output         datatypes.JSON `gorm:"type:jsonb"`
	Error          string         `gorm:"type:text"`
	RetryCount     int            `gorm:"default:0"`
	StartedAt      *time.Time
	CompletedAt    *time.Time
	database.AuditModel

	Execution DagExecution `gorm:"foreignKey:DagExecutionID"`
}
