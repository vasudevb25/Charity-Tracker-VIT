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

// type NGOHandler struct {
// 	service *services.NGOService
// }

// func NewNGOHandler(service *services.NGOService) *NGOHandler {
// 	return &NGOHandler{service: service}
// }

// func (h *NGOHandler) CreateNGO(c *gin.Context) {
// 	var ngo models.NGO
// 	if err := c.ShouldBindJSON(&ngo); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.CreateNGO(ctx, &ngo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "NGO created successfully", "ngo_id": ngo.ID.Hex()})
// }

func (h *NGOHandler) GetNGOs(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	ngos, err := h.service.GetNGOs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ngos)
}

func (h *NGOHandler) GetNGOByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	ngo, err := h.service.GetNGOByID(ctx, id)
	if err != nil {
		if err.Error() == "NGO not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ngo)
}

// func (h *NGOHandler) UpdateNGO(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	var ngo models.NGO
// 	if err := c.ShouldBindJSON(&ngo); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.UpdateNGO(ctx, id, &ngo); err != nil {
// 		if err.Error() == "NGO not found or no changes made" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "NGO updated successfully"})
// }

// func (h *NGOHandler) DeleteNGO(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.DeleteNGO(ctx, id); err != nil {
// 		if err.Error() == "NGO not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "NGO deleted successfully"})
// }

func (h *NGOHandler) SearchNGOs(c *gin.Context) {
	var req models.SearchNGOsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	ngos, err := h.service.SearchNGOs(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(ngos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No NGOs found matching your criteria"})
		return
	}

	c.JSON(http.StatusOK, ngos)
}

type NGOHandler struct {
	service *services.NGOService
}

func NewNGOHandler(service *services.NGOService) *NGOHandler {
	return &NGOHandler{service: service}
}

func (h *NGOHandler) CreateNGO(c *gin.Context) {
	claims, err := middleware.GetUserClaims(c) // Get claims from context
	if err != nil || claims.Role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to create NGOs"})
		return
	}

	var ngo models.NGO
	if err := c.ShouldBindJSON(&ngo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateNGO(ctx, &ngo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "NGO created successfully", "ngo_id": ngo.ID.Hex()})
}

// ... other functions (GetNGOs, GetNGOByID, SearchNGOs) remain largely the same, no auth needed for them in public group

// UpdateNGO now requires NGORequired middleware, which checks if :id matches claims.NGOID
func (h *NGOHandler) UpdateNGO(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
		return
	}

	// Middleware NGORequired() already checked if claims.NGOID matches id.
	// If Admin, middleware ensures AdminRequired() or NGORequired() with Admin bypass passes.

	var ngo models.NGO
	if err := c.ShouldBindJSON(&ngo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IMPORTANT: Ensure the NGOID in the body (if provided) also matches the authenticated user's NGOID or the path ID.
	// This prevents an authenticated NGO user from accidentally updating *another* NGO's profile
	// if the frontend sends an incorrect NGOID in the body.
	if !ngo.ID.IsZero() && ngo.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NGO ID in body must match path ID"})
		return
	}
	ngo.ID = id // Ensure we're updating the correct document

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.UpdateNGO(ctx, id, &ngo); err != nil {
		if err.Error() == "NGO not found or no changes made" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NGO updated successfully"})
}

// DeleteNGO similar to UpdateNGO, middleware handles initial check.
func (h *NGOHandler) DeleteNGO(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
		return
	}

	// Middleware NGORequired() already checked if claims.NGOID matches id (if NGO role).
	// If Admin, middleware ensures AdminRequired() or NGORequired() with Admin bypass passes.

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.DeleteNGO(ctx, id); err != nil {
		if err.Error() == "NGO not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "NGO deleted successfully"})
}
