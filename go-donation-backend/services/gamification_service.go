package services

import (
	"context"
	"fmt"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/models"
	"go-donation-backend/utils"
)

type GamificationService struct {
	achievementsCollection *mongo.Collection
	donationsCollection    *mongo.Collection // To calculate leaderboard
}

func NewGamificationService(client *mongo.Client) *GamificationService {
	return &GamificationService{
		achievementsCollection: config.GetCollection(client, "achievements"),
		donationsCollection:    config.GetCollection(client, "donations"),
	}
}

func (s *GamificationService) AwardDonationAchievement(ctx context.Context, donorID string, ngoID primitive.ObjectID, amount float64) {
	// Simplified: Award a badge for any donation, could be more complex (e.g., amount tiers, specific causes)
	badgeName := "Generous Donor"
	if amount > 1000 {
		badgeName = "Philanthropy Champion"
	}

	achievement := models.Achievement{
		ID:              utils.GenerateObjectID(),
		DonorID:         donorID,
		BadgeName:       badgeName,
		Cause:           "", // Can derive from NGO category later
		AchievementDate: utils.GetCurrentTime(),
	}

	_, err := s.achievementsCollection.InsertOne(ctx, achievement)
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to award achievement for donor %s", donorID))
	}
}

func (s *GamificationService) GetDonorAchievements(ctx context.Context, donorID string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	cursor, err := s.achievementsCollection.Find(ctx, bson.M{"donor_id": donorID})
	if err != nil {
		utils.LogError(err, fmt.Sprintf("Failed to get achievements for donor %s", donorID))
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &achievements); err != nil {
		utils.LogError(err, "Failed to decode donor achievements")
		return nil, err
	}
	return achievements, nil
}

func (s *GamificationService) GetLeaderboard(ctx context.Context, limit int64) ([]models.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Aggregate total donation amount per donor
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":          "$donor_id",
				"total_amount": bson.M{"$sum": "$amount"},
			},
		},
		{
			"$sort": bson.M{"total_amount": -1}, // Sort by total amount descending
		},
		{
			"$limit": limit,
		},
	}

	cursor, err := s.donationsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.LogError(err, "Failed to aggregate leaderboard data")
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaderboard []models.LeaderboardEntry
	if err = cursor.All(ctx, &leaderboard); err != nil {
		utils.LogError(err, "Failed to decode leaderboard entries")
		return nil, err
	}

	// Ensure stable sort if amounts are equal (optional, but good practice)
	sort.SliceStable(leaderboard, func(i, j int) bool {
		return leaderboard[i].TotalAmount > leaderboard[j].TotalAmount
	})

	return leaderboard, nil
}
