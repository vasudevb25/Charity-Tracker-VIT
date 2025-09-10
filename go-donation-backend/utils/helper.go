package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateObjectID creates a new MongoDB ObjectID.
func GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// GetCurrentTime returns the current UTC time.
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// LogError logs an error with a timestamp.
func LogError(err error, msg string) {
	log.Printf("[%s] ERROR: %s - %v", GetCurrentTime().Format(time.RFC3339), msg, err)
}

// ContextWithTimeout creates a context with a default timeout.
func ContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
