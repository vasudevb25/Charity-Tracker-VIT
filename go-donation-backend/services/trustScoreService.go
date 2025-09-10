package services

import (
	"go-donation-backend/models"
)

// CalculateTrustScore calculates the overall trust score
func CalculateTrustScore(ts models.TrustScore) float64 {
	// Weights for each factor (normalize them to percentages)
	weights := map[string]float64{
		"InflowOutflow":      0.10,
		"AuditScore":         0.50,
		"ReviewScore":        0.20,
		"Beneficiaries":      0.10,
		"SocialMediaEngage":  0.10,
		"TimeInOperation":    0.05,
		"ComplianceScore":    0.05,
		"DonationFrequency":  0.05,
	}

	// Calculate weighted sum for each factor
	score := (ts.InflowOutflow * weights["InflowOutflow"]) +
		(ts.AuditScore * weights["AuditScore"]) +
		(ts.ReviewScore * weights["ReviewScore"]) +
		(ts.Beneficiaries * weights["Beneficiaries"]) +
		(ts.SocialMediaEngage * weights["SocialMediaEngage"]) +
		(ts.TimeInOperation * weights["TimeInOperation"]) +
		(ts.ComplianceScore * weights["ComplianceScore"]) +
		(ts.DonationFrequency * weights["DonationFrequency"])

	return score
}
