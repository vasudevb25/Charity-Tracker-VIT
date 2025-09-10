package handlers

import (
	"github.com/gin-gonic/gin"
	"go-donation-backend/services"
	"net/http"
	"go-donation-backend/models"
)

// TrustScoreHandler handles the POST request to calculate trust score
func TrustScoreHandler(c *gin.Context) {
	var ts models.TrustScore


	if err := c.ShouldBindJSON(&ts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	score := services.CalculateTrustScore(ts)

	c.JSON(http.StatusOK, gin.H{"trust_score": score})
}
