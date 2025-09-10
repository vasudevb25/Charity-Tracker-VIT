package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/models"
	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type ExpenditureHandler struct {
	service *services.ExpenditureService
}

func NewExpenditureHandler(service *services.ExpenditureService) *ExpenditureHandler {
	return &ExpenditureHandler{service: service}
}

// func (h *ExpenditureHandler) AddExpenditure(c *gin.Context) {
// 	ngoIDParam := c.Param("ngoID")
// 	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	var expenditure models.Expenditure
// 	if err := c.ShouldBindJSON(&expenditure); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	expenditure.NGOID = ngoID // Ensure NGOID from path is used

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	if err := h.service.AddExpenditure(ctx, &expenditure); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Expenditure added successfully", "expenditure_id": expenditure.ID.Hex()})
// }

// func (h *ExpenditureHandler) GetExpendituresByNGO(c *gin.Context) {
// 	ngoIDParam := c.Param("ngoID")
// 	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format"})
// 		return
// 	}

// 	ctx, cancel := utils.ContextWithTimeout()
// 	defer cancel()

// 	expenditures, err := h.service.GetExpendituresByNGO(ctx, ngoID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if len(expenditures) == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "No expenditures found for this NGO"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, expenditures)
// }

func (h *ExpenditureHandler) AddExpenditure(c *gin.Context) {
	ngoIDParam := c.Param("ngoID")
	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format in path"})
		return
	}

	var expenditure models.Expenditure
	if err := c.ShouldBindJSON(&expenditure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Double-check: If NGOID is provided in the body, it must match the path parameter.
	// This prevents an authenticated NGO user from accidentally logging an expenditure
	// for a different NGO by putting a different NGOID in the JSON body.
	if !expenditure.NGOID.IsZero() && expenditure.NGOID != ngoID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: NGO ID in request body must match NGO ID in path"})
		return
	}
	expenditure.NGOID = ngoID // Ensure the expenditure is linked to the NGO from the path

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.AddExpenditure(ctx, &expenditure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expenditure added successfully", "expenditure_id": expenditure.ID.Hex()})
}

// GetExpendituresByNGO retrieves all expenditure records for a specific NGO.
// This route is protected by middleware.NGORequired() in main.go,
// which handles authorization.
func (h *ExpenditureHandler) GetExpendituresByNGO(c *gin.Context) {
	ngoIDParam := c.Param("ngoID")
	ngoID, err := primitive.ObjectIDFromHex(ngoIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format in path"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	expenditures, err := h.service.GetExpendituresByNGO(ctx, ngoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(expenditures) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No expenditures found for this NGO"})
		return
	}

	c.JSON(http.StatusOK, expenditures)
}
