package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-donation-backend/models"
	"go-donation-backend/services"
	"go-donation-backend/utils"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT upon successful registration
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), user.NGOID, user.DonorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User registered successfully",
		"user_id":  user.ID.Hex(),
		"role":     user.Role,
		"ngo_id":   user.NGOID.Hex(),
		"donor_id": user.DonorID,
		"token":    token,
	})
}

// Login handles user login and issues a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := utils.ContextWithTimeout()
	defer cancel()

	user, err := h.userService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role), user.NGOID, user.DonorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"user_id":  user.ID.Hex(),
		"role":     user.Role,
		"ngo_id":   user.NGOID.Hex(),
		"donor_id": user.DonorID,
		"token":    token,
	})
}
