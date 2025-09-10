package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Donation represents a single donation transaction
type Donation struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DonorID            string             `bson:"donor_id" json:"donor_id" binding:"required"` // Placeholder for user ID
	NGOID              primitive.ObjectID `bson:"ngo_id" json:"ngo_id" binding:"required"`
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
	RecipientType   string               `json:"recipient_type" binding:"required"` // "single_ngo", "multiple_ngos", "emergency_fund"
	NGOID           primitive.ObjectID   `json:"ngo_id,omitempty"`                  // For single NGO donation
	NGOIDs          []primitive.ObjectID `json:"ngo_ids,omitempty"`                 // For multiple NGO donation
	SplitAmounts    map[string]float64   `json:"split_amounts,omitempty"`           // For multiple NGOs, NGO_ID -> Amount
	EmergencyFundID primitive.ObjectID   `json:"emergency_fund_id,omitempty"`       // For emergency fund
}
