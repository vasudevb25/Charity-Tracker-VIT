package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"
)

type OrganizationService struct {
	collection *mongo.Collection
}

func NewOrganizationService(client *mongo.Client) *OrganizationService {
	return &OrganizationService{
		collection: config.GetCollection(client, "organizations"), // <--- Updated collection name
	}
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, organization *models.Organization) error {
	organization.ID = utils.GenerateObjectID()
	organization.CreatedAt = utils.GetCurrentTime()
	organization.UpdatedAt = utils.GetCurrentTime()
	organization.IsVerified = false // Default
	organization.TrustScore = 0.0   // Initialize TrustScore

	_, err := s.collection.InsertOne(ctx, organization)
	if err != nil {
		utils.LogError(err, "Failed to create organization")
		return err
	}
	return nil
}

func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id primitive.ObjectID) (*models.Organization, error) {
	var organization models.Organization
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&organization)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("organization not found")
		}
		utils.LogError(err, "Failed to get organization by ID")
		return nil, err
	}
	return &organization, nil
}

func (s *OrganizationService) GetOrganizations(ctx context.Context) ([]models.Organization, error) {
	var organizations []models.Organization
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogError(err, "Failed to get all organizations")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &organizations); err != nil {
		utils.LogError(err, "Failed to decode organizations")
		return nil, err
	}
	return organizations, nil
}

func (s *OrganizationService) UpdateOrganization(ctx context.Context, id primitive.ObjectID, organization *models.Organization) error {
	organization.UpdatedAt = utils.GetCurrentTime()
	update := bson.M{
		"$set": bson.M{
			"name":          organization.Name,
			"type":          organization.Type, // <--- Added Type
			"description":   organization.Description,
			"category":      organization.Category,
			"contact_email": organization.ContactEmail,
			"website":       organization.Website,
			"address":       organization.Address,
			// "trust_score": organization.TrustScore, // Trust score might be updated by a separate process/admin
			"updated_at": organization.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, "Failed to update organization")
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("organization not found or no changes made")
	}
	return nil
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.LogError(err, "Failed to delete organization")
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("organization not found")
	}
	return nil
}

func (s *OrganizationService) SearchOrganizations(ctx context.Context, req *models.SearchOrganizationsRequest) ([]models.Organization, error) {
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Location != "" {
		filter["address"] = primitive.Regex{Pattern: req.Location, Options: "i"}
	}
	if req.Name != "" {
		filter["name"] = primitive.Regex{Pattern: req.Name, Options: "i"}
	}
	if req.Type != "" { // <--- Added Type filter
		filter["type"] = req.Type
	}

	var organizations []models.Organization
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		utils.LogError(err, "Failed to search organizations")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &organizations); err != nil {
		utils.LogError(err, "Failed to decode searched organizations")
		return nil, err
	}
	return organizations, nil
}
