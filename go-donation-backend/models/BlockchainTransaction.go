package models

type BlockchainTransaction struct {
	DonationID       string  `json:"donation_id"`                 // MongoDB donation ID
	DonorID          string  `json:"donor_id"`                    // Who donated
	OrganizationID   string  `json:"organization_id"`             // Recipient org ID
	Amount           float64 `json:"amount"`                      // Donation amount
	Currency         string  `json:"currency"`                    // Currency (e.g., INR, USD)
	DonationDate     string  `json:"donation_date"`               // ISO8601 string
	OriginalDonation string  `json:"original_donation,omitempty"` // For split donations
}
