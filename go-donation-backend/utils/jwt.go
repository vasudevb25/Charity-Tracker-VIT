package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Claims defines the JWT claims for a user
type Claims struct {
	UserID         primitive.ObjectID `json:"user_id"`
	Email          string             `json:"email"`
	Role           string             `json:"role"`
	OrganizationID primitive.ObjectID `json:"organization_id,omitempty"` // <--- Updated to OrganizationID
	DonorID        string             `json:"donor_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for the given user claims
func GenerateToken(userID primitive.ObjectID, email string, role string, organizationID primitive.ObjectID, donorID string) (string, error) { // <--- Updated parameter
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretjwtkey" // Fallback for dev, DO NOT USE IN PROD
		LogError(fmt.Errorf("JWT_SECRET not set, using default. This is insecure for production."), "JWT Warning")
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID:         userID,
		Email:          email,
		Role:           role,
		OrganizationID: organizationID, // <--- Updated field
		DonorID:        donorID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		LogError(err, "Failed to sign JWT token")
		return "", err
	}
	return tokenString, nil
}

// ValidateToken parses and validates a JWT token, returning the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretjwtkey" // Fallback for dev, DO NOT USE IN PROD
		LogError(fmt.Errorf("JWT_SECRET not set, using default. This is insecure for production."), "JWT Warning")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		LogError(err, "Failed to parse or validate JWT token")
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// HashPassword hashes a plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		LogError(err, "Failed to hash password")
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain-text password with a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		LogError(err, "Password hash comparison failed") // Could be wrong password
	}
	return err == nil
}
