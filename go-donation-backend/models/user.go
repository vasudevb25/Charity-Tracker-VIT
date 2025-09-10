// ### Step 3: Create `models/user.go`

// This file defines the `User` struct, user roles, and request bodies for registration and login.

// ```go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole defines the type for user roles
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleNGO   UserRole = "ngo"
	RoleDonor UserRole = "donor"
)

// User represents a user of the platform
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string             `bson:"email" json:"email" binding:"required,email"`
	PasswordHash string             `bson:"password_hash" json:"-"` // Stored hashed password, never expose in JSON
	Role         UserRole           `bson:"role" json:"role" binding:"required"`
	NGOID        primitive.ObjectID `bson:"ngo_id,omitempty" json:"ngo_id,omitempty"`     // Link to NGO profile if role is NGO
	DonorID      string             `bson:"donor_id,omitempty" json:"donor_id,omitempty"` // Link to Donor profile (if different from User.ID)
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// RegisterRequest for user registration
type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role" binding:"required"`
	NGOID    string   `json:"ngo_id,omitempty"`   // String to allow binding, will convert to ObjectID
	DonorID  string   `json:"donor_id,omitempty"` // String
}

// LoginRequest for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
