package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmergencyFund represents a fund for crisis situations
type EmergencyFund struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name" binding:"required"`
	Description     string             `bson:"description" json:"description"`
	TargetAmount    float64            `bson:"target_amount,omitempty" json:"target_amount,omitempty"`
	CollectedAmount float64            `bson:"collected_amount" json:"collected_amount"`
	Status          string             `bson:"status" json:"status"` // e.g., "active", "closed", "met_target"
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	ClosedAt        *time.Time         `bson:"closed_at,omitempty" json:"closed_at,omitempty"` // Pointer for optional closing date
}
