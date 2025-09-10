package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type CertificateHandler struct {
	service *services.CertificateService
}

func NewCertificateHandler(service *services.CertificateService) *CertificateHandler {
	return &CertificateHandler{service: service}
}

func (h *CertificateHandler) GetDonationCertificate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	certificate, err := h.service.GetDonationCertificateByID(ctx, id)
	if err != nil {
		if err.Error() == "donation certificate not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certificate)
}

func (h *CertificateHandler) GetCertificateByDonationID(c *gin.Context) {
	donationIDParam := c.Param("donationID")
	donationID, err := primitive.ObjectIDFromHex(donationIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	certificate, err := h.service.GetCertificateByDonationID(ctx, donationID)
	if err != nil {
		if err.Error() == "certificate for donation not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, certificate)
}
