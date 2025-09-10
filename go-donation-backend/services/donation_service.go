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

type DonationService struct {
	collection           *mongo.Collection
	ngoService           *NGOService
	emergencyFundService *EmergencyFundService
	certificateService   *CertificateService
	gamificationService  *GamificationService
}

func NewDonationService(client *mongo.Client) *DonationService {
	return &DonationService{
		collection:           config.GetCollection(client, "donations"),
		ngoService:           NewNGOService(client), // Dependency injection
		emergencyFundService: NewEmergencyFundService(client),
		certificateService:   NewCertificateService(client),
		gamificationService:  NewGamificationService(client),
	}
}

func (s *DonationService) CreateDonation(ctx context.Context, req *models.DonationRequest) ([]models.Donation, error) {
	var createdDonations []models.Donation
	donationDate := utils.GetCurrentTime()

	switch req.RecipientType {
	case "single_ngo":
		if req.NGOID.IsZero() {
			return nil, errors.New("NGO ID is required for single NGO donation")
		}
		donation := models.Donation{
			ID:           utils.GenerateObjectID(),
			DonorID:      req.DonorID,
			NGOID:        req.NGOID,
			Amount:       req.Amount,
			Currency:     req.Currency,
			DonationDate: donationDate,
			Status:       "completed", // Simplified status
			IsSplit:      false,
			CreatedAt:    utils.GetCurrentTime(),
		}
		_, err := s.collection.InsertOne(ctx, donation)
		if err != nil {
			utils.LogError(err, "Failed to create single NGO donation")
			return nil, err
		}
		createdDonations = append(createdDonations, donation)
		s.certificateService.GenerateCertificate(ctx, &donation) // Generate certificate
		s.gamificationService.AwardDonationAchievement(ctx, req.DonorID, req.NGOID, req.Amount)

	case "multiple_ngos":
		if len(req.SplitAmounts) == 0 {
			return nil, errors.New("split amounts are required for multiple NGO donation")
		}
		originalDonationID := utils.GenerateObjectID()
		for ngoIDStr, amount := range req.SplitAmounts {
			ngoID, err := primitive.ObjectIDFromHex(ngoIDStr)
			if err != nil {
				utils.LogError(err, fmt.Sprintf("Invalid NGO ID in split amounts: %s", ngoIDStr))
				return nil, fmt.Errorf("invalid NGO ID: %s", ngoIDStr)
			}

			donation := models.Donation{
				ID:                 utils.GenerateObjectID(),
				DonorID:            req.DonorID,
				NGOID:              ngoID,
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
			s.gamificationService.AwardDonationAchievement(ctx, req.DonorID, ngoID, amount)
		}

	case "emergency_fund":
		if req.EmergencyFundID.IsZero() {
			return nil, errors.New("emergency fund ID is required for emergency fund donation")
		}
		// First, create the donation record
		donation := models.Donation{
			ID:           utils.GenerateObjectID(),
			DonorID:      req.DonorID,
			NGOID:        req.EmergencyFundID, // Using EmergencyFundID as NGOID for tracking
			Amount:       req.Amount,
			Currency:     req.Currency,
			DonationDate: donationDate,
			Status:       "completed",
			IsSplit:      false,
			CreatedAt:    utils.GetCurrentTime(),
		}
		_, err := s.collection.InsertOne(ctx, donation)
		if err != nil {
			utils.LogError(err, "Failed to create emergency fund donation")
			return nil, err
		}
		createdDonations = append(createdDonations, donation)
		s.gamificationService.AwardDonationAchievement(ctx, req.DonorID, req.EmergencyFundID, req.Amount)

		// Then, update the emergency fund's collected amount
		err = s.emergencyFundService.AddAmountToEmergencyFund(ctx, req.EmergencyFundID, req.Amount)
		if err != nil {
			// Handle potential consistency issues, maybe log and proceed, or try to revert donation
			utils.LogError(err, fmt.Sprintf("Failed to update collected amount for emergency fund %s after donation", req.EmergencyFundID.Hex()))
			// Consider adding a compensation/rollback mechanism here in a real application
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

func (s *DonationService) GetDonationsByNGO(ctx context.Context, ngoID primitive.ObjectID) ([]models.Donation, error) {
	var donations []models.Donation
	cursor, err := s.collection.Find(ctx, bson.M{"ngo_id": ngoID})
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get donations for NGO %s", ngoID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &donations); err != nil {
		utils.LogError(err, "Failed to decode donations for NGO")
		return nil, err
	}
	return donations, nil
}
