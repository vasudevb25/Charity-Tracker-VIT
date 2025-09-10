// package handlers

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson/primitive"

// 	"go-donation-backend/models"
// 	"go-donation-backend/services"
// 	"go-donation-backend/utils"
// )

// type ProjectUpdateHandler struct {
// 	service *services.ProjectUpdateService
// }

// func NewProjectUpdateHandler(service *services.ProjectUpdateService) *ProjectUpdateHandler {
// 	return &ProjectUpdateHandler{service: service}
// }

// func (h *ProjectUpdateHandler) CreateProjectUpdate(c *gin.Context) {
// 	ngoIDParam := c.Param("ngoID")
// 	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	var update models.ProjectUpdate
// 	if err := c.ShouldBindJSON(&update); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	update.NGOID = ngoID // Ensure NGOID from path is used

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.CreateProjectUpdate(ctx, &update); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Project update created successfully", "update_id": update.ID.Hex()})
// }

// func (h *ProjectUpdateHandler) GetProjectUpdatesByNGO(c *gin.Context) {
// 	ngoIDParam := c.Param("ngoID")
// 	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	updates, err := h.service.GetProjectUpdatesByNGO(ctx, ngoID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if len(updates) == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "No project updates found for this NGO"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, updates)
// }

// func (h *ProjectUpdateHandler) UpdateProjectUpdate(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project Update ID format"})
// 		return
// 	}

// 	var update models.ProjectUpdate
// 	if err := c.ShouldBindJSON(&update); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.UpdateProjectUpdate(ctx, id, &update); err != nil {
// 		if err.Error() == "Project update not found or no changes made" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Project update updated successfully"})
// }

// func (h *ProjectUpdateHandler) DeleteProjectUpdate(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Project Update ID format"})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.DeleteProjectUpdate(ctx, id); err != nil {
// 		if err.Error() == "Project update not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{"message": "Project update deleted successfully"})
//	}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/middleware" // Import middleware
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

// CreateProjectUpdate handles the creation of a new project update by an NGO.
// This route is protected by middleware.NGORequired() in main.go,
// which ensures the authenticated user is an NGO/Admin and, if NGO,
// that claims.NGOID matches the :ngoID path parameter.
func (h *ProjectUpdateHandler) CreateProjectUpdate(c *gin.Context) {
	ngoIDParam := c.Param("ngoID")
	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format in path"})
		return
	}

	var update models.ProjectUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Double-check: If NGOID is provided in the body, it must match the path parameter.
	if !update.NGOID.IsZero() && update.NGOID != ngoID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: NGO ID in request body must match NGO ID in path"})
		return
	}
	update.NGOID = ngoID // Ensure the update is linked to the NGO from the path

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateProjectUpdate(ctx, &update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project update created successfully", "update_id": update.ID.Hex()})
}

// GetProjectUpdatesByNGO retrieves project updates for a specific NGO (publicly accessible).
func (h *ProjectUpdateHandler) GetProjectUpdatesByNGO(c *gin.Context) {
	ngoIDParam := c.Param("ngoID")
	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format in path"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	updates, err := h.service.GetProjectUpdatesByNGO(ctx, ngoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(updates) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No project updates found for this NGO"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()}) // Should not happen if AuthMiddleware is before
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

	// Authorization check: If the user is an NGO, their claims.NGOID must match the existing update's NGOID.
	// Admins bypass this specific check due to AdminRequired middleware's structure allowing Admins.
	if claims.Role == string(models.RoleNGO) && claims.NGOID != existingUpdate.NGOID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only update your own NGO's project updates"})
		return
	}

	var updateData models.ProjectUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent changing the NGOID of an existing project update through this endpoint.
	// If NGOID is provided in the body, it must match the original.
	if !updateData.NGOID.IsZero() && updateData.NGOID != existingUpdate.NGOID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot change NGO ID of an existing project update"})
		return
	}
	// Always ensure the NGOID is carried over from the existing record, even if not provided in updateData.
	updateData.NGOID = existingUpdate.NGOID

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

	// Authorization check: If the user is an NGO, their claims.NGOID must match the existing update's NGOID.
	if claims.Role == string(models.RoleNGO) && claims.NGOID != existingUpdate.NGOID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only delete your own NGO's project updates"})
		return
	}

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
