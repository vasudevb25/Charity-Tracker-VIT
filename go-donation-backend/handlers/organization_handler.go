package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/models"
	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type OrganizationHandler struct { // <--- Updated name
	service *services.OrganizationService // <--- Updated service
}

func NewOrganizationHandler(service *services.OrganizationService) *OrganizationHandler { // <--- Updated name
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) { // <--- Updated name
	// Authorization is handled by middleware.AdminRequired() in main.go
	var organization models.Organization // <--- Updated struct
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateOrganization(ctx, &organization); err != nil { // <--- Updated service method
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organization created successfully", "organization_id": organization.ID.Hex()}) // <--- Updated response field
}

func (h *OrganizationHandler) GetOrganizations(c *gin.Context) { // <--- Updated name
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	organizations, err := h.service.GetOrganizations(ctx) // <--- Updated service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, organizations)
}

func (h *OrganizationHandler) GetOrganizationByID(c *gin.Context) { // <--- Updated name
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	organization, err := h.service.GetOrganizationByID(ctx, id) // <--- Updated service method
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, organization)
}

func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) { // <--- Updated name
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format"})
		return
	}

	// Authorization is handled by middleware.OrganizationRequired() in main.go,
	// which checks if path :id matches claims.OrganizationID for organization role,
	// or if the user is an admin.

	var organization models.Organization // <--- Updated struct
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IMPORTANT: Ensure the OrganizationID in the body (if provided) also matches the path ID.
	// This prevents an authenticated Organization user from accidentally updating *another* organization's profile
	// if the frontend sends an incorrect OrganizationID in the body.
	if !organization.ID.IsZero() && organization.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID in body must match path ID"})
		return
	}
	organization.ID = id // Ensure we're updating the correct document

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.UpdateOrganization(ctx, id, &organization); err != nil { // <--- Updated service method
		if err.Error() == "organization not found or no changes made" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully"})
}

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) { // <--- Updated name
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format"})
		return
	}

	// Authorization is handled by middleware.OrganizationRequired() in main.go,
	// which checks if path :id matches claims.OrganizationID for organization role,
	// or if the user is an admin.

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.DeleteOrganization(ctx, id); err != nil { // <--- Updated service method
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

func (h *OrganizationHandler) SearchOrganizations(c *gin.Context) { // <--- Updated name
	var req models.SearchOrganizationsRequest // <--- Updated struct
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	organizations, err := h.service.SearchOrganizations(ctx, &req) // <--- Updated service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(organizations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No organizations found matching your criteria"})
		return
	}

	c.JSON(http.StatusOK, organizations)
}
