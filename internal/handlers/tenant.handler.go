package handlers

import (
	"backend-sevima/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/yuusufyan/go-common/response"
)

type TenantHandler struct {
	tenantService services.TenantService
}

// NewTenantHandler creates a new instance of TenantHandler
func NewTenantHandler(tenantService services.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// GetAll retrieves all tenants
// @Summary Get all tenants
// @Description Retrieve a list of all tenants in the system
// @Tags Tenant
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /tenants [get]
func (h *TenantHandler) GetAll(c *fiber.Ctx) error {
	tenants, err := h.tenantService.GetAll()
	if err != nil {
		return response.RespondWithError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Successfully retrieved tenants", tenants)
}

// GetByID retrieves a tenant by ID
// @Summary Get a tenant by ID
// @Description Retrieve a specific tenant's details
// @Tags Tenant
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 404 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	tenant, err := h.tenantService.GetByID(id)
	if err != nil {
		return response.RespondWithError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Successfully retrieved tenant", tenant)
}

// Create creates a new tenant
// @Summary Create a new tenant
// @Description Create a new tenant in the system
// @Tags Tenant
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body services.CreateTenantRequest true "Create Tenant Request"
// @Success 201 {object} response.Common
// @Failure 400 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /tenants [post]
func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var req services.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	tenant, err := h.tenantService.Create(req)
	if err != nil {
		return response.RespondWithError(c, err)
	}
	return response.Success(c, fiber.StatusCreated, "Tenant created successfully", tenant)
}

// Update updates an existing tenant
// @Summary Update an existing tenant
// @Description Update the details of an existing tenant
// @Tags Tenant
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Tenant ID"
// @Param request body services.UpdateTenantRequest true "Update Tenant Request"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 404 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /tenants/{id} [put]
func (h *TenantHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	tenant, err := h.tenantService.Update(id, req)
	if err != nil {
		return response.RespondWithError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Tenant updated successfully", tenant)
}

// Delete removes a tenant
// @Summary Delete a tenant
// @Description Remove a tenant from the system
// @Tags Tenant
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.Common
// @Failure 401 {object} response.Common
// @Failure 404 {object} response.Common
// @Failure 500 {object} response.Common
// @Router /tenants/{id} [delete]
func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.tenantService.Delete(id); err != nil {
		return response.RespondWithError(c, err)
	}
	return response.Success(c, fiber.StatusOK, "Tenant deleted successfully", nil)
}
