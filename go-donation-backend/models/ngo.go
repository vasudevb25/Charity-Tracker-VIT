package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NGO represents a Non-Governmental Organization
type NGO struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Description  string             `bson:"description" json:"description"`
	Category     []string           `bson:"category" json:"category"` // e.g., "Education", "Healthcare", "Environment"
	ContactEmail string             `bson:"contact_email" json:"contact_email" binding:"required,email"`
	Website      string             `bson:"website,omitempty" json:"website,omitempty"`
	Address      string             `bson:"address,omitempty" json:"address,omitempty"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"` // Verified NGO Onboarding (simplified)
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type SearchNGOsRequest struct {
	Category string `form:"category"`
	Location string `form:"location"`
	Name     string `form:"name"`
}
