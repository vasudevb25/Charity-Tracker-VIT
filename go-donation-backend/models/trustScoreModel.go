package models

type TrustScore struct {
	InflowOutflow      float64 // 10% weight
	AuditScore         float64 // 50% weight
	ReviewScore        float64 // 20% weight
	Beneficiaries      float64 // 10% weight
	SocialMediaEngage  float64 // 10% weight
	TimeInOperation    float64 // 5% weight
	ComplianceScore    float64 // 5% weight
	DonationFrequency  float64 // 5% weight
}
