package services

import (
	"context"
	"errors"
	"fmt"
	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportService struct {
	donationsCollection     *mongo.Collection
	expendituresCollection  *mongo.Collection
	organizationsCollection *mongo.Collection // To verify organization exists and get its name
}

func NewReportService(client *mongo.Client) *ReportService {
	return &ReportService{
		donationsCollection:     config.GetCollection(client, "donations"),
		expendituresCollection:  config.GetCollection(client, "expenditures"),
		organizationsCollection: config.GetCollection(client, "organizations"), // <--- Updated collection name
	}
}

// GetOrganizationAuditReport provides an aggregated financial report for an organization.
func (s *ReportService) GetOrganizationAuditReport(ctx context.Context, orgID primitive.ObjectID) (*models.OrganizationAuditReport, error) {
	// Verify organization exists and get its name
	var org models.Organization
	err := s.organizationsCollection.FindOne(ctx, bson.M{"_id": orgID}).Decode(&org)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("organization not found")
		}
		return nil, fmt.Errorf("failed to check organization existence: %w", err)
	}

	// Aggregate total donations for the organization
	donationsPipeline := []bson.M{
		{"$match": bson.M{"organization_id": orgID}}, // <--- Updated field
		{"$group": bson.M{
			"_id":          nil,
			"total_amount": bson.M{"$sum": "$amount"},
			"count":        bson.M{"$sum": 1},
		}},
	}
	var donationResults []struct {
		TotalAmount float64 `bson:"total_amount"`
		Count       int64   `bson:"count"`
	}
	cursor, err := s.donationsCollection.Aggregate(ctx, donationsPipeline)
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to aggregate donations for organization %s", orgID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &donationResults) // No error check needed, as it might be empty

	totalDonations := 0.0
	donationCount := int64(0)
	if len(donationResults) > 0 {
		totalDonations = donationResults[0].TotalAmount
		donationCount = donationResults[0].Count
	}

	// Aggregate total expenditures for the organization
	expendituresPipeline := []bson.M{
		{"$match": bson.M{"organization_id": orgID}}, // <--- Updated field
		{"$group": bson.M{
			"_id":          nil,
			"total_amount": bson.M{"$sum": "$amount"},
			"count":        bson.M{"$sum": 1},
		}},
	}
	var expenditureResults []struct {
		TotalAmount float64 `bson:"total_amount"`
		Count       int64   `bson:"count"`
	}
	cursor, err = s.expendituresCollection.Aggregate(ctx, expendituresPipeline)
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to aggregate expenditures for organization %s", orgID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &expenditureResults)

	totalExpenditures := 0.0
	expenditureCount := int64(0)
	if len(expenditureResults) > 0 {
		totalExpenditures = expenditureResults[0].TotalAmount
		expenditureCount = expenditureResults[0].Count
	}

	report := &models.OrganizationAuditReport{
		OrganizationID:          orgID,
		OrganizationName:        org.Name, // Added name
		TotalDonationsAmount:    totalDonations,
		DonationCount:           donationCount,
		TotalExpendituresAmount: totalExpenditures,
		ExpenditureCount:        expenditureCount,
		NetFunds:                totalDonations - totalExpenditures,
	}
	return report, nil
}
