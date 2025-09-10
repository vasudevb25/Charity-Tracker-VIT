package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collaboration represents a project run by multiple organizations
type Collaboration struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	OrganizationIDs []primitive.ObjectID `bson:"organization_ids" json:"organization_ids" binding:"required,min=2"` // <--- Updated to OrganizationIDs
	ProjectName     string               `bson:"project_name" json:"project_name" binding:"required"`
	Description     string               `bson:"description" json:"description"`
	StartDate       time.Time            `bson:"start_date" json:"start_date"`
	EndDate         time.Time            `bson:"end_date" json:"end_date"`
	Status          string               `bson:"status" json:"status"` // e.g., "planning", "active", "completed"
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
}
