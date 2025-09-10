package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-donation-backend/models"
	"go-donation-backend/utils"
)

// Gin context key for user claims
const UserClaimsKey = "userClaims"

// AuthMiddleware extracts and validates JWT from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context for subsequent handlers
		c.Set(UserClaimsKey, claims)
		c.Next()
	}
}

// AdminRequired checks if the authenticated user has the 'admin' role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(UserClaimsKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userClaims := claims.(*utils.Claims)
		if userClaims.Role != string(models.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// NGORequired checks if the authenticated user has the 'ngo' role
// It also checks if the 'ngoID' path parameter matches the user's NGOID.
func NGORequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(UserClaimsKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userClaims := claims.(*utils.Claims)
		if userClaims.Role != string(models.RoleNGO) && userClaims.Role != string(models.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "NGO or Admin access required"})
			c.Abort()
			return
		}

		// If an Admin is accessing an NGO-specific route, they can proceed without NGOID matching.
		// If an NGO user is accessing, their NGOID must match the one in the path.
		if userClaims.Role == string(models.RoleNGO) {
			paramNGOID := c.Param("ngoID")
			if paramNGOID == "" { // For routes that don't have :ngoID in path (e.g., /ngo/profile if implemented)
				c.Next()
				return
			}

			objID, err := primitive.ObjectIDFromHex(paramNGOID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NGO ID format in path"})
				c.Abort()
				return
			}

			if userClaims.NGOID != objID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only manage your own NGO resources"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// DonorRequired checks if the authenticated user has the 'donor' role
// It also checks if the 'donorID' path parameter matches the user's DonorID.
func DonorRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(UserClaimsKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userClaims := claims.(*utils.Claims)
		if userClaims.Role != string(models.RoleDonor) && userClaims.Role != string(models.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Donor or Admin access required"})
			c.Abort()
			return
		}

		// If an Admin is accessing a Donor-specific route, they can proceed without DonorID matching.
		// If a Donor user is accessing, their DonorID must match the one in the path.
		if userClaims.Role == string(models.RoleDonor) {
			paramDonorID := c.Param("donorID")
			if paramDonorID == "" {
				c.Next()
				return
			}

			if userClaims.DonorID != paramDonorID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You can only view your own donor resources"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// GetUserClaims retrieves claims from context (used inside handlers)
func GetUserClaims(c *gin.Context) (*utils.Claims, error) {
	claims, exists := c.Get(UserClaimsKey)
	if !exists {
		return nil, fmt.Errorf("user claims not found in context")
	}
	return claims.(*utils.Claims), nil
}
