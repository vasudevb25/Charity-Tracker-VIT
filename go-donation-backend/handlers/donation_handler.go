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

type DonationHandler struct {
	service *services.DonationService
}

func NewDonationHandler(service *services.DonationService) *DonationHandler {
	return &DonationHandler{service: service}
}

func (h *DonationHandler) CreateDonation(c *gin.Context) {
	var req models.DonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Authorization Check for DonorID ---
	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// If a Donor is making a donation, ensure req.DonorID matches their claims.DonorID
	// If an Admin is making a donation, they can specify any donor_id.
	if claims.Role == string(models.RoleDonor) && req.DonorID != claims.DonorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only make donations as yourself"})
		return
	}
	// If donor_id is empty in request, auto-fill it with claims.DonorID for convenience
	if req.DonorID == "" {
		req.DonorID = claims.DonorID
	}
	// --- End Authorization Check ---

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	// Create donation
	donations, err := h.service.CreateDonation(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call blockchain function to simulate sending data to blockchain
	transactionHash, err := services.SendToBlockchain(req.DonorID, req.OrganizationID.Hex(), req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send transaction to blockchain"})
		return
	}

	// Return success message with transaction hash
	c.JSON(http.StatusCreated, gin.H{
		"message":          "Donation(s) created successfully",
		"donations":        donations,
		"transaction_hash": transactionHash,
	})
}

// GetDonationsByDonor (already protected by middleware.DonorRequired() in main.go)
func (h *DonationHandler) GetDonationsByDonor(c *gin.Context) {
	donorID := c.Param("donorID")

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	donations, err := h.service.GetDonationsByDonor(ctx, donorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(donations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No donations found for this donor"})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetDonationsByOrganization(c *gin.Context) {
	orgIDParam := c.Param("orgID")
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	donations, err := h.service.GetDonationsByOrganization(ctx, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(donations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No donations found for this organization"})
		return
	}

	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetDonationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	donation, err := h.service.GetDonationByID(ctx, id)
	if err != nil {
		if err.Error() == "donation not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, donation)
}
