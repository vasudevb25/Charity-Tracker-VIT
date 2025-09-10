package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationType defines the type of institution
type OrganizationType string

const (
	TypeNGO         OrganizationType = "NGO"
	TypeOrphanage   OrganizationType = "Orphanage"
	TypeSchool      OrganizationType = "School"
	TypeHospital    OrganizationType = "Hospital"
	TypeTrust       OrganizationType = "Trust"
	TypeInstitution OrganizationType = "Institution"
	TypeOther       OrganizationType = "Other"
)

// Organization represents any type of institution (NGO, Orphanage, etc.)
type Organization struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Type         OrganizationType   `bson:"type" json:"type" binding:"required"` // Added Type
	Description  string             `bson:"description" json:"description"`
	Category     []string           `bson:"category" json:"category"` // e.g., "Education", "Healthcare", "Environment"
	ContactEmail string             `bson:"contact_email" json:"contact_email" binding:"required,email"`
	Website      string             `bson:"website,omitempty" json:"website,omitempty"`
	Address      string             `bson:"address,omitempty" json:"address,omitempty"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"`
	TrustScore   float64            `bson:"trust_score" json:"trust_score"` // Added Trust Score
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// SearchOrganizationsRequest for searching organizations
type SearchOrganizationsRequest struct {
	Category string `form:"category"`
	Location string `form:"location"`
	Name     string `form:"name"`
	Type     string `form:"type"` // Added Type filter
}
