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

type ReportHandler struct {
	reportService   *services.ReportService
	donationService *services.DonationService // For donor history
}

func NewReportHandler(reportService *services.ReportService, donationService *services.DonationService) *ReportHandler {
	return &ReportHandler{
		reportService:   reportService,
		donationService: donationService,
	}
}

// GetOrganizationAuditReport provides an aggregated financial report for an organization.
// This endpoint is protected and only accessible by the organization itself or an admin.
func (h *ReportHandler) GetOrganizationAuditReport(c *gin.Context) {
	orgIDParam := c.Param("orgID")
	orgID, err := primitive.ObjectIDFromHex(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Organization ID format"})
		return
	}

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Authorization check: If the user is an Organization, their claims.OrganizationID must match the orgID in the path.
	if claims.Role == string(models.RoleOrganization) && claims.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only view your own organization's audit report"})
		return
	}
	// Admin role is allowed by the middleware via OrganizationRequired.

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	report, err := h.reportService.GetOrganizationAuditReport(ctx, orgID)
	if err != nil {
		if err.Error() == "organization not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetDonorTransactionHistory provides a list of all donations made by a specific donor.
// This endpoint is protected and only accessible by the donor themselves or an admin.
func (h *ReportHandler) GetDonorTransactionHistory(c *gin.Context) {
	donorID := c.Param("donorID")

	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Authorization check: If the user is a Donor, their claims.DonorID must match the donorID in the path.
	if claims.Role == string(models.RoleDonor) && claims.DonorID != donorID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only view your own transaction history"})
		return
	}
	// Admin role is allowed by the middleware via DonorRequired.

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	donations, err := h.donationService.GetDonationsByDonor(ctx, donorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(donations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No transactions found for this donor"})
		return
	}

	c.JSON(http.StatusOK, donations)
}
