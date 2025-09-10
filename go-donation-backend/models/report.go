package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// OrganizationAuditReport represents an aggregated financial overview for an organization.
type OrganizationAuditReport struct {
	OrganizationID          primitive.ObjectID `json:"organization_id"`
	OrganizationName        string             `json:"organization_name"` // Added for display
	TotalDonationsAmount    float64            `json:"total_donations_amount"`
	DonationCount           int64              `json:"donation_count"`
	TotalExpendituresAmount float64            `json:"total_expenditures_amount"`
	ExpenditureCount        int64              `json:"expenditure_count"`
	NetFunds                float64            `json:"net_funds"`
	// Additional fields can be added here, e.g., breakdowns by category, top donors, etc.
}
