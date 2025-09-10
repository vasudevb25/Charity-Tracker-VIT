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

type ProjectUpdateService struct {
	collection *mongo.Collection
}

func NewProjectUpdateService(client *mongo.Client) *ProjectUpdateService {
	return &ProjectUpdateService{
		collection: config.GetCollection(client, "project_updates"),
	}
}

func (s *ProjectUpdateService) CreateProjectUpdate(ctx context.Context, update *models.ProjectUpdate) error {
	update.ID = utils.GenerateObjectID()
	update.CreatedAt = utils.GetCurrentTime()
	update.UpdateDate = utils.GetCurrentTime() // Or taken from request
	update.UpdatedAt = utils.GetCurrentTime()

	_, err := s.collection.InsertOne(ctx, update)
	if err != nil {
		utils.LogError(err, "Failed to create project update")
		return err
	}
	return nil
}

func (s *ProjectUpdateService) GetProjectUpdatesByOrganization(ctx context.Context, organizationID primitive.ObjectID) ([]models.ProjectUpdate, error) { // <--- Updated function name
	var updates []models.ProjectUpdate
	cursor, err := s.collection.Find(ctx, bson.M{"organization_id": organizationID}) // <--- Updated field
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get project updates for organization %s", organizationID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &updates); err != nil {
		utils.LogError(err, "Failed to decode project updates for organization")
		return nil, err
	}
	return updates, nil
}

func (s *ProjectUpdateService) UpdateProjectUpdate(ctx context.Context, id primitive.ObjectID, updateData *models.ProjectUpdate) error {
	updateData.UpdatedAt = utils.GetCurrentTime()
	update := bson.M{
		"$set": bson.M{
			"title":           updateData.Title,
			"description":     updateData.Description,
			"media_urls":      updateData.MediaURLs,
			"update_date":     updateData.UpdateDate,
			"updated_at":      updateData.UpdatedAt,
			"organization_id": updateData.OrganizationID, // Ensure OrganizationID is updated if provided (but handler prevents change)
		},
	}

	result, err := s.collection.UpdateByID(ctx, id, update)
	if err != nil {
		utils.LogError(err, "Failed to update project update")
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Project update not found or no changes made")
	}
	return nil
}

func (s *ProjectUpdateService) DeleteProjectUpdate(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		utils.LogError(err, "Failed to delete project update")
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Project update not found")
	}
	return nil
}

// GetProjectUpdateByID retrieves a project update by its ID.
func (s *ProjectUpdateService) GetProjectUpdateByID(ctx context.Context, id primitive.ObjectID) (*models.ProjectUpdate, error) {
	var update models.ProjectUpdate
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("project update not found")
		}
		utils.LogError(err, fmt.Sprintf("Failed to get project update by ID %s", id.Hex()))
		return nil, err
	}
	return &update, nil
}
