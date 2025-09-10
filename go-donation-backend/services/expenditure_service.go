package services

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"
)

type ExpenditureService struct {
	collection *mongo.Collection
}

func NewExpenditureService(client *mongo.Client) *ExpenditureService {
	return &ExpenditureService{
		collection: config.GetCollection(client, "expenditures"),
	}
}

func (s *ExpenditureService) AddExpenditure(ctx context.Context, expenditure *models.Expenditure) error {
	expenditure.ID = utils.GenerateObjectID()
	expenditure.CreatedAt = utils.GetCurrentTime()
	expenditure.ExpenditureDate = utils.GetCurrentTime() // Or taken from request

	_, err := s.collection.InsertOne(ctx, expenditure)
	if err != nil {
		utils.LogError(err, "Failed to add expenditure")
		return err
	}
	return nil
}

func (s *ExpenditureService) GetExpendituresByOrganization(ctx context.Context, organizationID primitive.ObjectID) ([]models.Expenditure, error) { // <--- Updated function name
	var expenditures []models.Expenditure
	cursor, err := s.collection.Find(ctx, bson.M{"organization_id": organizationID}) // <--- Updated field
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get expenditures for organization %s", organizationID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &expenditures); err != nil {
		utils.LogError(err, "Failed to decode expenditures for organization")
		return nil, err
	}
	return expenditures, nil
}

func (s *ExpenditureService) GetExpendituresByDonation(ctx context.Context, donationID primitive.ObjectID) ([]models.Expenditure, error) {
	var expenditures []models.Expenditure
	cursor, err := s.collection.Find(ctx, bson.M{"donation_id": donationID})
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get expenditures for donation %s", donationID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &expenditures); err != nil {
		utils.LogError(err, "Failed to decode expenditures for donation")
		return nil, err
	}
	return expenditures, nil
}
