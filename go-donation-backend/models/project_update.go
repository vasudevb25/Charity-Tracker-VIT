package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectUpdate represents an update from an organization about project progress
type ProjectUpdate struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organization_id" json:"organization_id" binding:"required"` // <--- Updated to OrganizationID
	Title          string             `bson:"title" json:"title" binding:"required"`
	Description    string             `bson:"description" json:"description"`
	MediaURLs      []string           `bson:"media_urls,omitempty" json:"media_urls,omitempty"` // Links to images/videos (IPFS concept)
	UpdateDate     time.Time          `bson:"update_date" json:"update_date"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
