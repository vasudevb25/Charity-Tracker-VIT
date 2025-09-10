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

type UserService struct {
	collection *mongo.Collection
	ngoService *NGOService // Dependency to check/link NGO existence
}

func NewUserService(client *mongo.Client, ngoService *NGOService) *UserService {
	return &UserService{
		collection: config.GetCollection(client, "users"),
		ngoService: ngoService,
	}
}

// CreateUser registers a new user
func (s *UserService) CreateUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Check if user with this email already exists
	existingUser, err := s.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err // Database error
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           utils.GenerateObjectID(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		CreatedAt:    utils.GetCurrentTime(),
		UpdatedAt:    utils.GetCurrentTime(),
	}

	// Logic for linking NGOID or DonorID based on role
	if req.Role == models.RoleNGO {
		if req.NGOID == "" {
			return nil, errors.New("ngo_id is required for NGO role")
		}
		ngoObjID, err := primitive.ObjectIDFromHex(req.NGOID)
		if err != nil {
			return nil, errors.New("invalid ngo_id format")
		}
		// Optionally, verify NGOID exists in the NGOs collection
		_, err = s.ngoService.GetNGOByID(ctx, ngoObjID)
		if err != nil {
			return nil, fmt.Errorf("invalid ngo_id provided: %w", err)
		}
		user.NGOID = ngoObjID
	} else if req.Role == models.RoleDonor {
		// For simplicity, DonorID can be the UserID or a separate profile ID.
		// Here, we'll link it directly to the user's ID for simplicity,
		// but a more complex system might have a separate Donor profile.
		if req.DonorID != "" {
			user.DonorID = req.DonorID
		} else {
			user.DonorID = user.ID.Hex() // Default to user's ID as donor_id
		}
	} else if req.Role == models.RoleAdmin {
		// No special linking for Admin, but ensure no NGOID/DonorID is mistakenly set
		user.NGOID = primitive.NilObjectID
		user.DonorID = ""
	}

	_, err = s.collection.InsertOne(ctx, user)
	if err != nil {
		utils.LogError(err, "Failed to create user")
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err // mongo.ErrNoDocuments if not found
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err // mongo.ErrNoDocuments if not found
	}
	return &user, nil
}

// AuthenticateUser verifies credentials and returns the user if successful
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid credentials") // User not found
		}
		utils.LogError(err, "Error retrieving user for authentication")
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials") // Password mismatch
	}

	return user, nil
}
