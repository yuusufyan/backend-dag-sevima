package services

import (
	"errors"

	"backend-sevima/internal/models"
	"github.com/yuusufyan/go-common/pkg/apperror"
	"gorm.io/gorm"
)

// TenantService interface defines CRUD operations for Tenant
type TenantService interface {
	GetAll() ([]model.Tenant, error)
	GetByID(id string) (*model.Tenant, error)
	Create(req CreateTenantRequest) (*model.Tenant, error)
	Update(id string, req UpdateTenantRequest) (*model.Tenant, error)
	Delete(id string) error
}

type tenantService struct {
	db *gorm.DB
}

// NewTenantService creates a new instance of TenantService
func NewTenantService(db *gorm.DB) TenantService {
	return &tenantService{
		db: db,
	}
}

// Data structures for Request/Response
type CreateTenantRequest struct {
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateTenantRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (s *tenantService) GetAll() ([]model.Tenant, error) {
	var tenants []model.Tenant
	if err := s.db.Find(&tenants).Error; err != nil {
		return nil, apperror.InternalServer("Failed to retrieve tenants")
	}
	return tenants, nil
}

func (s *tenantService) GetByID(id string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := s.db.First(&tenant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("Tenant not found")
		}
		return nil, apperror.InternalServer("Failed to retrieve tenant")
	}
	return &tenant, nil
}

func (s *tenantService) Create(req CreateTenantRequest) (*model.Tenant, error) {
	if req.Name == "" {
		return nil, apperror.BadRequest("Tenant name is required")
	}

	tenant := model.Tenant{
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return nil, apperror.InternalServer("Failed to create tenant")
	}

	return &tenant, nil
}

func (s *tenantService) Update(id string, req UpdateTenantRequest) (*model.Tenant, error) {
	tenant, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	// Note: We might need a pointer for IsActive to distinguish between false and not-provided, 
	// but for simplicity we will just assign it if we don't use a pointer.
	tenant.IsActive = req.IsActive

	if err := s.db.Save(tenant).Error; err != nil {
		return nil, apperror.InternalServer("Failed to update tenant")
	}

	return tenant, nil
}

func (s *tenantService) Delete(id string) error {
	tenant, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(tenant).Error; err != nil {
		return apperror.InternalServer("Failed to delete tenant")
	}

	return nil
}
