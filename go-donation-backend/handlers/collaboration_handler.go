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

type CollaborationHandler struct {
	service *services.CollaborationService
}

func NewCollaborationHandler(service *services.CollaborationService) *CollaborationHandler {
	return &CollaborationHandler{service: service}
}

// Helper function to check if an ObjectID is present in a slice of ObjectIDs
func containsObjectID(slice []primitive.ObjectID, val primitive.ObjectID) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// CreateCollaboration handles the creation of a new collaboration.
// This route is protected by middleware.AuthMiddleware() and is nested under authenticatedGroup.Group("/organization").
// Therefore, middleware.OrganizationRequired() will run for this route, ensuring the user is Organization/Admin.
// Further authorization logic for Organization participants is handled within this method.
func (h *CollaborationHandler) CreateCollaboration(c *gin.Context) {
	var collab models.Collaboration
	if err := c.ShouldBindJSON(&collab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must be one of the participating Organization IDs in the request body.
	if claims.Role == string(models.RoleOrganization) { // <--- Updated role
		if !containsObjectID(collab.OrganizationIDs, claims.OrganizationID) { // <--- Updated field
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: As an Organization, you must be a participant in the collaboration you create"})
			return
		}
	}
	// If Admin, they can create a collaboration with any Organization IDs.

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateCollaboration(ctx, &collab); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Collaboration created successfully", "collaboration_id": collab.ID.Hex()})
}

// GetCollaborations retrieves all collaborations (publicly accessible).
func (h *CollaborationHandler) GetCollaborations(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	collabs, err := h.service.GetCollaborations(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(collabs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No collaborations found"})
		return
	}

	c.JSON(http.StatusOK, collabs)
}

// GetCollaborationsByOrganization retrieves collaborations a specific Organization is part of.
// This route is protected by middleware.OrganizationRequired() in main.go,
// which handles authorization.
func (h *CollaborationHandler) GetCollaborationsByOrganization(c *gin.Context) { // <--- Updated function name
	orgIDParam := c.Param("orgID") // <--- Updated path parameter
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format in path"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	collabs, err := h.service.GetCollaborationsByOrganization(ctx, orgID) // <--- Updated service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(collabs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No collaborations found for this organization"})
		return
	}

	c.JSON(http.StatusOK, collabs)
}

// UpdateCollaboration handles updating an existing collaboration.
// This route is protected by middleware.AuthMiddleware() in main.go.
// Further authorization is handled within this method.
func (h *CollaborationHandler) UpdateCollaboration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collaboration ID format in path"})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	// Retrieve the existing collaboration to verify ownership/participation
	existingCollab, err := h.service.GetCollaborationByID(ctx, id)
	if err != nil {
		if err.Error() == "collaboration not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must be one of the participating Organization IDs in the existing collaboration.
	if claims.Role == string(models.RoleOrganization) { // <--- Updated role
		if !containsObjectID(existingCollab.OrganizationIDs, claims.OrganizationID) { // <--- Updated field
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You are not a participant in this collaboration"})
			return
		}
	}
	// If Admin, they can update any collaboration.

	var collabData models.Collaboration
	if err := c.ShouldBindJSON(&collabData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional check: If an Organization is updating, they cannot remove themselves from the OrganizationIDs list.
	if claims.Role == string(models.RoleOrganization) { // <--- Updated role
		if !containsObjectID(collabData.OrganizationIDs, claims.OrganizationID) { // <--- Updated field
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You cannot remove your Organization from a collaboration you are updating"})
			return
		}
	}

	// Prevent changing the ID of an existing collaboration.
	if !collabData.ID.IsZero() && collabData.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Collaboration ID in body must match path ID"})
		return
	}
	collabData.ID = id // Ensure we're updating the correct document

	if err := h.service.UpdateCollaboration(ctx, id, &collabData); err != nil {
		if err.Error() == "Collaboration not found or no changes made" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaboration updated successfully"})
}

// DeleteCollaboration handles deleting an existing collaboration.
// This route is protected by middleware.AuthMiddleware() in main.go.
// Further authorization is handled within this method.
func (h *CollaborationHandler) DeleteCollaboration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Collaboration ID format in path"})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	// Retrieve the existing collaboration to verify ownership/participation
	existingCollab, err := h.service.GetCollaborationByID(ctx, id)
	if err != nil {
		if err.Error() == "collaboration not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must be one of the participating Organization IDs in the existing collaboration.
	if claims.Role == string(models.RoleOrganization) { // <--- Updated role
		if !containsObjectID(existingCollab.OrganizationIDs, claims.OrganizationID) { // <--- Updated field
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You are not a participant in this collaboration"})
			return
		}
	}
	// If Admin, they can delete any collaboration.

	if err := h.service.DeleteCollaboration(ctx, id); err != nil {
		if err.Error() == "Collaboration not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaboration deleted successfully"})
}
