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

// AddExpenditure handles the creation of a new expenditure record by an Organization.
// This route is protected by middleware.OrganizationRequired() in main.go,
// which ensures the authenticated user is an Organization/Admin and, if Organization,
// that claims.OrganizationID matches the :orgID path parameter.
func (h *ExpenditureHandler) AddExpenditure(c *gin.Context) {
	orgIDParam := c.Param("id") // <--- Updated path parameter
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format in path"})
		return
	}

	var expenditure models.Expenditure
	if err := c.ShouldBindJSON(&expenditure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Double-check: If OrganizationID is provided in the body, it must match the path parameter.
	if !expenditure.OrganizationID.IsZero() && expenditure.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Organization ID in request body must match Organization ID in path"})
		return
	}
	expenditure.OrganizationID = orgID // Ensure the expenditure is linked to the Organization from the path

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.AddExpenditure(ctx, &expenditure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expenditure added successfully", "expenditure_id": expenditure.ID.Hex()})
}

// GetExpendituresByOrganization retrieves all expenditure records for a specific Organization.
// This route is protected by middleware.OrganizationRequired() in main.go,
// which handles authorization.
func (h *ExpenditureHandler) GetExpendituresByOrganization(c *gin.Context) { // <--- Updated function name
	orgIDParam := c.Param("orgID") // <--- Updated path parameter
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format in path"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	expenditures, err := h.service.GetExpendituresByOrganization(ctx, orgID) // <--- Updated service method
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(expenditures) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No expenditures found for this organization"})
		return
	}

	c.JSON(http.StatusOK, expenditures)
}
