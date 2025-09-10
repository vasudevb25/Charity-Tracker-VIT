package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Expenditure represents how an organization spends donated money
type Expenditure struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrganizationID  primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"` // <--- Updated to OrganizationID
	DonationID      primitive.ObjectID `bson:"donation_id,omitempty" json:"donation_id,omitempty"`        // Optional: Link to a specific donation
	Amount          float64            `bson:"amount" json:"amount" binding:"required,gte=0"`
	Description     string             `bson:"description" json:"description"`
	ProofURL        string             `bson:"proof_url,omitempty" json:"proof_url,omitempty"` // Link to receipt/document (IPFS concept)
	ExpenditureDate time.Time          `bson:"expenditure_date" json:"expenditure_date"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}
