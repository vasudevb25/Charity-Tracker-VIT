package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DonationCertificate represents a proof of impact NFT (simplified)
type DonationCertificate struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DonationID     primitive.ObjectID `bson:"donation_id" json:"donation_id" binding:"required"`
	DonorID        string             `bson:"donor_id" json:"donor_id" binding:"required"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"` // <--- Updated to OrganizationID
	CertificateURL string             `bson:"certificate_url" json:"certificate_url"`                    // URL to a generated digital certificate
	IssueDate      time.Time          `bson:"issue_date" json:"issue_date"`
	Metadata       map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"` // e.g., "ImpactStatement": "Funded 5 meals"
}
