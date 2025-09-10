package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type GamificationHandler struct {
	service *services.GamificationService
}

func NewGamificationHandler(service *services.GamificationService) *GamificationHandler {
	return &GamificationHandler{service: service}
}

func (h *GamificationHandler) GetDonorAchievements(c *gin.Context) {
	donorID := c.Param("donorID")

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	achievements, err := h.service.GetDonorAchievements(ctx, donorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(achievements) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No achievements found for this donor"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

func (h *GamificationHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit if not provided or invalid
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	leaderboard, err := h.service.GetLeaderboard(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(leaderboard) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Leaderboard is empty"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}
