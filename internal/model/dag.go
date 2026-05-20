package model

import (
	"time"

	"gorm.io/datatypes"
)

// DagTemplate stores the definition of a DAG.
type DagTemplate struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID    string         `gorm:"index;not null"`
	Name        string         `gorm:"size:255;not null"`
	Description string         `gorm:"type:text"`
	// Definition contains the raw JSON structure of the DAG (nodes and edges)
	Definition  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DagExecution tracks a single execution instance of a DagTemplate.
type DagExecution struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID      string         `gorm:"index;not null"`
	DagTemplateID string         `gorm:"type:uuid;index;not null"`
	Status        string         `gorm:"size:50;not null"` // PENDING, RUNNING, SUCCESS, FAILED
	Input         datatypes.JSON `gorm:"type:jsonb"`
	Output        datatypes.JSON `gorm:"type:jsonb"`
	Error         string         `gorm:"type:text"`
	StartedAt     *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Template DagTemplate `gorm:"foreignKey:DagTemplateID"`
}

// TaskInstance tracks the execution of a single Task within a DagExecution.
type TaskInstance struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
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
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Execution DagExecution `gorm:"foreignKey:DagExecutionID"`
}
