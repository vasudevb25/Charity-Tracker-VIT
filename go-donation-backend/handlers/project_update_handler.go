package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/middleware"
	"go-donation-backend/models"
	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type ProjectUpdateHandler struct {
	service *services.ProjectUpdateService
}

func NewProjectUpdateHandler(service *services.ProjectUpdateService) *ProjectUpdateHandler {
	return &ProjectUpdateHandler{service: service}
}

// CreateProjectUpdate handles the creation of a new project update by an Organization.
// This route is protected by middleware.OrganizationRequired() in main.go,
// which ensures the authenticated user is an Organization/Admin and, if Organization,
// that claims.OrganizationID matches the :orgID path parameter.
func (h *ProjectUpdateHandler) CreateProjectUpdate(c *gin.Context) {
	orgIDParam := c.Param("orgID") // <--- Updated path parameter
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format in path"})
		return
	}

	var update models.ProjectUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Double-check: If OrganizationID is provided in the body, it must match the path parameter.
	if !update.OrganizationID.IsZero() && update.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Organization ID in request body must match Organization ID in path"})
		return
	}
	update.OrganizationID = orgID // Ensure the update is linked to the Organization from the path

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateProjectUpdate(ctx, &update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project update created successfully", "update_id": update.ID.Hex()})
}

// GetProjectUpdatesByOrganization retrieves project updates for a specific Organization (publicly accessible).
func (h *ProjectUpdateHandler) GetProjectUpdatesByOrganization(c *gin.Context) { // <--- Updated function name
	orgIDParam := c.Param("orgID") // <--- Updated path parameter
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format in path"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	updates, err := h.service.GetProjectUpdatesByOrganization(ctx, orgID) // <--- Updated service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(updates) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No project updates found for this organization"})
		return
	}

	c.JSON(http.StatusOK, updates)
}

// UpdateProjectUpdate handles updating an existing project update.
// This route is protected by middleware.AuthMiddleware() in main.go.
// Further authorization is handled within this method.
func (h *ProjectUpdateHandler) UpdateProjectUpdate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project Update ID format in path"})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	// Retrieve the existing project update to verify ownership
	existingUpdate, err := h.service.GetProjectUpdateByID(ctx, id)
	if err != nil {
		if err.Error() == "project update not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must match the existing update's OrganizationID.
	if claims.Role == string(models.RoleOrganization) && claims.OrganizationID != existingUpdate.OrganizationID { // <--- Updated role and field
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only update your own organization's project updates"})
		return
	}
	// Admin role is allowed by the middleware.

	var updateData models.ProjectUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent changing the OrganizationID of an existing project update through this endpoint.
	// If OrganizationID is provided in the body, it must match the original.
	if !updateData.OrganizationID.IsZero() && updateData.OrganizationID != existingUpdate.OrganizationID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot change Organization ID of an existing project update"})
		return
	}
	// Always ensure the OrganizationID is carried over from the existing record, even if not provided in updateData.
	updateData.OrganizationID = existingUpdate.OrganizationID

	if err := h.service.UpdateProjectUpdate(ctx, id, &updateData); err != nil {
		if err.Error() == "Project update not found or no changes made" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project update updated successfully"})
}

// DeleteProjectUpdate handles deleting an existing project update.
// This route is protected by middleware.AuthMiddleware() in main.go.
// Further authorization is handled within this method.
func (h *ProjectUpdateHandler) DeleteProjectUpdate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project Update ID format in path"})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	// Retrieve the existing project update to verify ownership
	existingUpdate, err := h.service.GetProjectUpdateByID(ctx, id)
	if err != nil {
		if err.Error() == "project update not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must match the existing update's OrganizationID.
	if claims.Role == string(models.RoleOrganization) && claims.OrganizationID != existingUpdate.OrganizationID { // <--- Updated role and field
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only delete your own organization's project updates"})
		return
	}
	// Admin role is allowed by the middleware.

	if err := h.service.DeleteProjectUpdate(ctx, id); err != nil {
		if err.Error() == "Project update not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project update deleted successfully"})
}
