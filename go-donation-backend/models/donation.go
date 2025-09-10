package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Donation represents a single donation transaction
type Donation struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DonorID            string             `bson:"donor_id" json:"donor_id" binding:"required"`               // Placeholder for user ID
	OrganizationID     primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"` // <--- Updated to OrganizationID
	Amount             float64            `bson:"amount" json:"amount" binding:"required,gte=0"`
	Currency           string             `bson:"currency" json:"currency" binding:"required"`
	DonationDate       time.Time          `bson:"donation_date" json:"donation_date"`
	Status             string             `bson:"status" json:"status"` // e.g., "completed", "pending", "failed"
	IsSplit            bool               `bson:"is_split" json:"is_split"`
	OriginalDonationID primitive.ObjectID `bson:"original_donation_id,omitempty" json:"original_donation_id,omitempty"` // For split parts
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
}

// DonationRequest for accepting donations from donors
type DonationRequest struct {
	DonorID         string               `json:"donor_id" binding:"required"`
	Amount          float64              `json:"amount" binding:"required,gte=0"`
	Currency        string               `json:"currency" binding:"required"`
	RecipientType   string               `json:"recipient_type" binding:"required"` // "single_organization", "multiple_organizations"
	OrganizationID  primitive.ObjectID   `json:"organization_id,omitempty"`         // For single organization donation
	OrganizationIDs []primitive.ObjectID `json:"organization_ids,omitempty"`        // For multiple organizations donation
	SplitAmounts    map[string]float64   `json:"split_amounts,omitempty"`           // For multiple organizations, Org_ID -> Amount
}
