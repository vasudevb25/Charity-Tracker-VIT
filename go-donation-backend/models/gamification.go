package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Achievement represents a donor's badge/achievement
type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DonorID         string             `bson:"donor_id" json:"donor_id" binding:"required"`
	BadgeName       string             `bson:"badge_name" json:"badge_name" binding:"required"`
	Cause           string             `bson:"cause,omitempty" json:"cause,omitempty"` // e.g., "Education Hero"
	AchievementDate time.Time          `bson:"achievement_date" json:"achievement_date"`
}

// LeaderboardEntry represents a single entry in a leaderboard
type LeaderboardEntry struct {
	DonorID     string  `bson:"_id" json:"donor_id"` // Group by donor_id
	TotalAmount float64 `bson:"total_amount" json:"total_amount"`
}
