package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/models"
	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type EmergencyFundHandler struct {
	service *services.EmergencyFundService
}

func NewEmergencyFundHandler(service *services.EmergencyFundService) *EmergencyFundHandler {
	return &EmergencyFundHandler{service: service}
}

func (h *EmergencyFundHandler) CreateEmergencyFund(c *gin.Context) {
	var fund models.EmergencyFund
	if err := c.ShouldBindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.CreateEmergencyFund(ctx, &fund); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Emergency fund created successfully", "fund_id": fund.ID.Hex()})
}

func (h *EmergencyFundHandler) GetEmergencyFunds(c *gin.Context) {
	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	funds, err := h.service.GetEmergencyFunds(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(funds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No emergency funds found"})
		return
	}

	c.JSON(http.StatusOK, funds)
}

func (h *EmergencyFundHandler) GetEmergencyFundByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fund ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	fund, err := h.service.GetEmergencyFundByID(ctx, id)
	if err != nil {
		if err.Error() == "emergency fund not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fund)
}

func (h *EmergencyFundHandler) UpdateEmergencyFund(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Emergency Fund ID format"})
		return
	}

	var fund models.EmergencyFund
	if err := c.ShouldBindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.UpdateEmergencyFund(ctx, id, &fund); err != nil {
		if err.Error() == "Emergency fund not found or no changes made" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emergency fund updated successfully"})
}

func (h *EmergencyFundHandler) DeleteEmergencyFund(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Emergency Fund ID format"})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	if err := h.service.DeleteEmergencyFund(ctx, id); err != nil {
		if err.Error() == "Emergency fund not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emergency fund deleted successfully"})
}
