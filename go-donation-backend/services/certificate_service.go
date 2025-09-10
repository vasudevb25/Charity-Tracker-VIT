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

type CertificateService struct {
	collection *mongo.Collection
}

func NewCertificateService(client *mongo.Client) *CertificateService {
	return &CertificateService{
		collection: config.GetCollection(client, "donation_certificates"),
	}
}

// GenerateCertificate creates a simplified proof of impact (no actual NFT minting)
func (s *CertificateService) GenerateCertificate(ctx context.Context, donation *models.Donation) (*models.DonationCertificate, error) {
	certificate := models.DonationCertificate{
		ID:             utils.GenerateObjectID(),
		DonationID:     donation.ID,
		DonorID:        donation.DonorID,
		NGOID:          donation.NGOID,
		CertificateURL: fmt.Sprintf("https://example.com/certificates/%s", utils.GenerateObjectID().Hex()), // Placeholder URL
		IssueDate:      utils.GetCurrentTime(),
		Metadata: map[string]string{
			"ImpactStatement": fmt.Sprintf("Donation of %.2f %s contributed to cause by NGO %s", donation.Amount, donation.Currency, donation.NGOID.Hex()),
		},
	}

	_, err := s.collection.InsertOne(ctx, certificate)
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to generate certificate for donation %s", donation.ID.Hex()))
		return nil, err
	}
	return &certificate, nil
}

func (s *CertificateService) GetDonationCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.DonationCertificate, error) {
	var certificate models.DonationCertificate
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&certificate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("donation certificate not found")
		}
		utils.LogError(err, "Failed to get donation certificate by ID")
		return nil, err
	}
	return &certificate, nil
}

func (s *CertificateService) GetCertificateByDonationID(ctx context.Context, donationID primitive.ObjectID) (*models.DonationCertificate, error) {
	var certificate models.DonationCertificate
	err := s.collection.FindOne(ctx, bson.M{"donation_id": donationID}).Decode(&certificate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("certificate for donation not found")
		}
		utils.LogError(err, "Failed to get certificate by donation ID")
		return nil, err
	}
	return &certificate, nil
}
