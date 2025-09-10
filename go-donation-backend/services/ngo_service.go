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

type NGOService struct {
	collection *mongo.Collection
}

func NewNGOService(client *mongo.Client) *NGOService {
	return &NGOService{
		collection: config.GetCollection(client, "ngos"),
	}
}

func (s *NGOService) CreateNGO(ctx context.Context, ngo *models.NGO) error {
	ngo.ID = utils.GenerateObjectID()
	ngo.CreatedAt = utils.GetCurrentTime()
	ngo.UpdatedAt = utils.GetCurrentTime()
	// Simplified: IsVerified could be false by default and set to true via an admin process
	ngo.IsVerified = false

	_, err := s.collection.InsertOne(ctx, ngo)
	if err != nil {
		utils.LogError(err, "Failed to create NGO")
		return err
	}
	return nil
}

func (s *NGOService) GetNGOByID(ctx context.Context, id primitive.ObjectID) (*models.NGO, error) {
	var ngo models.NGO
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ngo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("NGO not found")
		}
		utils.LogError(err, "Failed to get NGO by ID")
		return nil, err
	}
	return &ngo, nil
}

func (s *NGOService) GetNGOs(ctx context.Context) ([]models.NGO, error) {
	var ngos []models.NGO
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogError(err, "Failed to get all NGOs")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &ngos); err != nil {
		utils.LogError(err, "Failed to decode NGOs")
		return nil, err
	}
	return ngos, nil
}

func (s *NGOService) UpdateNGO(ctx context.Context, id primitive.ObjectID, ngo *models.NGO) error {
	ngo.UpdatedAt = utils.GetCurrentTime()
	update := bson.M{
		"$set": bson.M{
			"name":          ngo.Name,
			"description":   ngo.Description,
			"category":      ngo.Category,
			"contact_email": ngo.ContactEmail,
			"website":       ngo.Website,
			"address":       ngo.Address,
			"updated_at":    ngo.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, "Failed to update NGO")
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("NGO not found or no changes made")
	}
	return nil
}

func (s *NGOService) DeleteNGO(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.LogError(err, "Failed to delete NGO")
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("NGO not found")
	}
	return nil
}

func (s *NGOService) SearchNGOs(ctx context.Context, req *models.SearchNGOsRequest) ([]models.NGO, error) {
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Location != "" {
		// This assumes 'Address' field contains searchable location data
		filter["address"] = primitive.Regex{Pattern: req.Location, Options: "i"}
	}
	if req.Name != "" {
		filter["name"] = primitive.Regex{Pattern: req.Name, Options: "i"}
	}

	var ngos []models.NGO
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		utils.LogError(err, "Failed to search NGOs")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &ngos); err != nil {
		utils.LogError(err, "Failed to decode searched NGOs")
		return nil, err
	}
	return ngos, nil
}
