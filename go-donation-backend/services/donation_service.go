package services

import (
	"context"
	"errors"
	"fmt"

	// Re-added log
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"
)

type DonationService struct {
	collection          *mongo.Collection
	organizationService *OrganizationService // <--- Dependency injection
	certificateService  *CertificateService
	gamificationService *GamificationService
}

func NewDonationService(client *mongo.Client) *DonationService {
	return &DonationService{
		collection:          config.GetCollection(client, "donations"),
		organizationService: NewOrganizationService(client), // Dependency injection
		certificateService:  NewCertificateService(client),
		gamificationService: NewGamificationService(client),
	}
}

func (s *DonationService) CreateDonation(ctx context.Context, req *models.DonationRequest) ([]models.Donation, error) {
	var createdDonations []models.Donation
	donationDate := utils.GetCurrentTime()

	switch req.RecipientType {
	case "single_organization": // <--- Updated recipient type
		if req.OrganizationID.IsZero() {
			return nil, errors.New("organization ID is required for single organization donation")
		}
		donation := models.Donation{
			ID:             utils.GenerateObjectID(),
			DonorID:        req.DonorID,
			OrganizationID: req.OrganizationID, // <--- Updated field
			Amount:         req.Amount,
			Currency:       req.Currency,
			DonationDate:   donationDate,
			Status:         "completed", // Simplified status
			IsSplit:        false,
			CreatedAt:      utils.GetCurrentTime(),
		}
		_, err := s.collection.InsertOne(ctx, donation)
		if err != nil {
			utils.LogError(err, "Failed to create single organization donation")
			return nil, err
		}
		createdDonations = append(createdDonations, donation)
		s.certificateService.GenerateCertificate(ctx, &donation) // Generate certificate
		s.gamificationService.AwardDonationAchievement(ctx, req.DonorID, req.OrganizationID, req.Amount)

	case "multiple_organizations": // <--- Updated recipient type
		if len(req.SplitAmounts) == 0 {
			return nil, errors.New("split amounts are required for multiple organization donation")
		}
		originalDonationID := utils.GenerateObjectID()
		for orgIDStr, amount := range req.SplitAmounts {
			orgID, err := primitive.ObjectIDFromHex(orgIDStr)
			if err != nil {
				utils.LogError(err, fmt.Sprintf("Invalid Organization ID in split amounts: %s", orgIDStr))
				return nil, fmt.Errorf("invalid Organization ID: %s", orgIDStr)
			}

			donation := models.Donation{
				ID:                 utils.GenerateObjectID(),
				DonorID:            req.DonorID,
				OrganizationID:     orgID, // <--- Updated field
				Amount:             amount,
				Currency:           req.Currency,
				DonationDate:       donationDate,
				Status:             "completed",
				IsSplit:            true,
				OriginalDonationID: originalDonationID,
				CreatedAt:          utils.GetCurrentTime(),
			}
			_, err = s.collection.InsertOne(ctx, donation)
			if err != nil {
				utils.LogError(err, "Failed to create split donation part")
				return nil, err
			}
			createdDonations = append(createdDonations, donation)
			s.certificateService.GenerateCertificate(ctx, &donation) // Generate certificate for each split
			s.gamificationService.AwardDonationAchievement(ctx, req.DonorID, orgID, amount)
		}

	default:
		return nil, errors.New("invalid recipient type")
	}

	return createdDonations, nil
}

func (s *DonationService) GetDonationByID(ctx context.Context, id primitive.ObjectID) (*models.Donation, error) {
	var donation models.Donation
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&donation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("donation not found")
		}
		utils.LogError(err, "Failed to get donation by ID")
		return nil, err
	}
	return &donation, nil
}

func (s *DonationService) GetDonationsByDonor(ctx context.Context, donorID string) ([]models.Donation, error) {
	var donations []models.Donation
	cursor, err := s.collection.Find(ctx, bson.M{"donor_id": donorID})
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get donations for donor %s", donorID))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &donations); err != nil {
		utils.LogError(err, "Failed to decode donations for donor")
		return nil, err
	}
	return donations, nil
}

func (s *DonationService) GetDonationsByOrganization(ctx context.Context, organizationID primitive.ObjectID) ([]models.Donation, error) { // <--- Updated function name
	var donations []models.Donation
	cursor, err := s.collection.Find(ctx, bson.M{"organization_id": organizationID}) // <--- Updated field
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get donations for organization %s", organizationID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &donations); err != nil {
		utils.LogError(err, "Failed to decode donations for organization")
		return nil, err
	}
	return donations, nil
}
