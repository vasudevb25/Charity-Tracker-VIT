package services

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"
)

type EmergencyFundService struct {
	collection *mongo.Collection
}

func NewEmergencyFundService(client *mongo.Client) *EmergencyFundService {
	return &EmergencyFundService{
		collection: config.GetCollection(client, "emergency_funds"),
	}
}

func (s *EmergencyFundService) CreateEmergencyFund(ctx context.Context, fund *models.EmergencyFund) error {
	fund.ID = utils.GenerateObjectID()
	fund.CreatedAt = utils.GetCurrentTime()
	fund.CollectedAmount = 0.0 // Initialize
	fund.Status = "active"     // Default status

	_, err := s.collection.InsertOne(ctx, fund)
	if err != nil {
		utils.LogError(err, "Failed to create emergency fund")
		return err
	}
	return nil
}

func (s *EmergencyFundService) GetEmergencyFundByID(ctx context.Context, id primitive.ObjectID) (*models.EmergencyFund, error) {
	var fund models.EmergencyFund
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&fund)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("emergency fund not found")
		}
		utils.LogError(err, "Failed to get emergency fund by ID")
		return nil, err
	}
	return &fund, nil
}

func (s *EmergencyFundService) GetEmergencyFunds(ctx context.Context) ([]models.EmergencyFund, error) {
	var funds []models.EmergencyFund
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogError(err, "Failed to get all emergency funds")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &funds); err != nil {
		utils.LogError(err, "Failed to decode emergency funds")
		return nil, err
	}
	return funds, nil
}

// AddAmountToEmergencyFund updates the collected amount of an emergency fund.
func (s *EmergencyFundService) AddAmountToEmergencyFund(ctx context.Context, id primitive.ObjectID, amount float64) error {
	update := bson.M{
		"$inc": bson.M{"collected_amount": amount},
	}
	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to add amount to emergency fund %s", id.Hex()))
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Emergency fund not found or no changes made")
	}
	return nil
}

func (s *EmergencyFundService) UpdateEmergencyFund(ctx context.Context, id primitive.ObjectID, fund *models.EmergencyFund) error {
	updateFields := bson.M{
		"name":          fund.Name,
		"description":   fund.Description,
		"target_amount": fund.TargetAmount,
		"status":        fund.Status,
	}

	if fund.Status == "closed" && fund.ClosedAt == nil {
		now := utils.GetCurrentTime()
		fund.ClosedAt = &now
		updateFields["closed_at"] = fund.ClosedAt
	} else if fund.Status != "closed" {
		updateFields["closed_at"] = nil // Clear closed_at if fund is reopened
	}

	update := bson.M{"$set": updateFields}

	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, "Failed to update emergency fund")
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Emergency fund not found or no changes made")
	}
	return nil
}

func (s *EmergencyFundService) DeleteEmergencyFund(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.LogError(err, "Failed to delete emergency fund")
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Emergency fund not found")
	}
	return nil
}
