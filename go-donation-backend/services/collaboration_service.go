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

type CollaborationService struct {
	collection *mongo.Collection
}

func NewCollaborationService(client *mongo.Client) *CollaborationService {
	return &CollaborationService{
		collection: config.GetCollection(client, "collaborations"),
	}
}

func (s *CollaborationService) CreateCollaboration(ctx context.Context, collab *models.Collaboration) error {
	collab.ID = utils.GenerateObjectID()
	collab.CreatedAt = utils.GetCurrentTime()
	collab.UpdatedAt = utils.GetCurrentTime()
	collab.Status = "planning" // Default status

	_, err := s.collection.InsertOne(ctx, collab)
	if err != nil {
		utils.LogError(err, "Failed to create collaboration")
		return err
	}
	return nil
}

func (s *CollaborationService) GetCollaborationByID(ctx context.Context, id primitive.ObjectID) (*models.Collaboration, error) {
	var collab models.Collaboration
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&collab)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("collaboration not found")
		}
		utils.LogError(err, fmt.Sprintf("Failed to get collaboration by ID %s", id.Hex()))
		return nil, err
	}
	return &collab, nil
}

func (s *CollaborationService) GetCollaborations(ctx context.Context) ([]models.Collaboration, error) {
	var collabs []models.Collaboration
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogError(err, "Failed to get all collaborations")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &collabs); err != nil {
		utils.LogError(err, "Failed to decode collaborations")
		return nil, err
	}
	return collabs, nil
}

func (s *CollaborationService) GetCollaborationsByOrganization(ctx context.Context, organizationID primitive.ObjectID) ([]models.Collaboration, error) { // <--- Updated function name
	var collabs []models.Collaboration
	cursor, err := s.collection.Find(ctx, bson.M{"organization_ids": organizationID}) // <--- Updated field
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get collaborations for organization %s", organizationID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &collabs); err != nil {
		utils.LogError(err, "Failed to decode collaborations for organization")
		return nil, err
	}
	return collabs, nil
}

func (s *CollaborationService) UpdateCollaboration(ctx context.Context, id primitive.ObjectID, collabData *models.Collaboration) error {
	collabData.UpdatedAt = utils.GetCurrentTime()
	update := bson.M{
		"$set": bson.M{
			"organization_ids": collabData.OrganizationIDs, // <--- Updated field
			"project_name":     collabData.ProjectName,
			"description":      collabData.Description,
			"start_date":       collabData.StartDate,
			"end_date":         collabData.EndDate,
			"status":           collabData.Status,
			"updated_at":       collabData.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, "Failed to update collaboration")
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Collaboration not found or no changes made")
	}
	return nil
}

func (s *CollaborationService) DeleteCollaboration(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.LogError(err, "Failed to delete collaboration")
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Collaboration not found")
	}
	return nil
}
