// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/mongo"

// 	"go-donation-backend/config"
// 	"go-donation-backend/handlers"
// 	"go-donation-backend/services"
// )

// var mongoClient *mongo.Client

// func main() {
// 	// Load environment variables
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using environment variables directly.")
// 	}

// 	// Connect to MongoDB
// 	var err error
// 	mongoClient, err = config.ConnectDB()
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}
// 	defer func() {
// 		if err = mongoClient.Disconnect(context.Background()); err != nil {
// 			log.Fatalf("Error disconnecting from MongoDB: %v", err)
// 		}
// 	}()

// 	// Ping the DB to verify connection
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = mongoClient.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to ping MongoDB: %v", err)
// 	}
// 	log.Println("Connected to MongoDB!")

// 	// Initialize services with the MongoDB client
// 	ngoService := services.NewNGOService(mongoClient)
// 	donationService := services.NewDonationService(mongoClient)
// 	expenditureService := services.NewExpenditureService(mongoClient)
// 	projectUpdateService := services.NewProjectUpdateService(mongoClient)
// 	certificateService := services.NewCertificateService(mongoClient)
// 	gamificationService := services.NewGamificationService(mongoClient)
// 	collaborationService := services.NewCollaborationService(mongoClient)
// 	emergencyFundService := services.NewEmergencyFundService(mongoClient)

// 	// Initialize handlers with their respective services
// 	ngoHandler := handlers.NewNGOHandler(ngoService)
// 	donationHandler := handlers.NewDonationHandler(donationService)
// 	expenditureHandler := handlers.NewExpenditureHandler(expenditureService)
// 	projectUpdateHandler := handlers.NewProjectUpdateHandler(projectUpdateService)
// 	certificateHandler := handlers.NewCertificateHandler(certificateService)
// 	gamificationHandler := handlers.NewGamificationHandler(gamificationService)
// 	collaborationHandler := handlers.NewCollaborationHandler(collaborationService)
// 	emergencyFundHandler := handlers.NewEmergencyFundHandler(emergencyFundService)

// 	// Setup Gin router
// 	router := gin.Default()
// 	router.Use(cors.New(cors.Config{
// 		// Allow requests from your local HTML server's origin
// 		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"}, // Adjust if your local server uses a different port/hostname
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // Add any custom headers your frontend might send
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour, // Cache preflight requests for 12 hours
// 	}))

// 	// Public Routes (e.g., for donors to view)
// 	public := router.Group("/api/v1")
// 	{
// 		// NGOs
// 		public.GET("/ngos", ngoHandler.GetNGOs)
// 		public.GET("/ngos/:id", ngoHandler.GetNGOByID)
// 		public.GET("/ngos/search", ngoHandler.SearchNGOs) // AI Chatbot simplified

// 		// Donations
// 		public.GET("/donations/donor/:donorID", donationHandler.GetDonationsByDonor)
// 		public.GET("/donations/ngo/:ngoID", donationHandler.GetDonationsByNGO)
// 		public.GET("/donations/:id", donationHandler.GetDonationByID)

// 		// Project Updates
// 		public.GET("/project-updates/ngo/:ngoID", projectUpdateHandler.GetProjectUpdatesByNGO)

// 		// Certificates
// 		public.GET("/certificates/:id", certificateHandler.GetDonationCertificate)
// 		public.GET("/certificates/donation/:donationID", certificateHandler.GetCertificateByDonationID)

// 		// Gamification
// 		public.GET("/gamification/donor/:donorID/achievements", gamificationHandler.GetDonorAchievements)
// 		public.GET("/gamification/leaderboard", gamificationHandler.GetLeaderboard)

// 		// Collaborations
// 		public.GET("/collaborations", collaborationHandler.GetCollaborations)
// 		public.GET("/collaborations/ngo/:ngoID", collaborationHandler.GetCollaborationsByNGO)

// 		// Emergency Funds
// 		public.GET("/emergency-funds", emergencyFundHandler.GetEmergencyFunds)
// 		public.GET("/emergency-funds/:id", emergencyFundHandler.GetEmergencyFundByID)
// 	}

// 	// NGO-specific Routes (e.g., authenticated NGO users)
// 	ngoAuth := router.Group("/api/v1/ngo")
// 	{
// 		// NGOs (for NGO to manage its own profile)
// 		ngoAuth.POST("/register", ngoHandler.CreateNGO) // Simplified NGO Onboarding
// 		ngoAuth.PUT("/:id", ngoHandler.UpdateNGO)
// 		ngoAuth.DELETE("/:id", ngoHandler.DeleteNGO)

// 		// Expenditures
// 		ngoAuth.POST("/:ngoID/expenditures", expenditureHandler.AddExpenditure)
// 		ngoAuth.GET("/:ngoID/expenditures", expenditureHandler.GetExpendituresByNGO)

// 		// Project Updates
// 		ngoAuth.POST("/:ngoID/project-updates", projectUpdateHandler.CreateProjectUpdate)
// 		ngoAuth.PUT("/project-updates/:id", projectUpdateHandler.UpdateProjectUpdate)
// 		ngoAuth.DELETE("/project-updates/:id", projectUpdateHandler.DeleteProjectUpdate)

// 		// Collaborations
// 		ngoAuth.POST("/collaborations", collaborationHandler.CreateCollaboration)
// 		ngoAuth.PUT("/collaborations/:id", collaborationHandler.UpdateCollaboration)
// 		ngoAuth.DELETE("/collaborations/:id", collaborationHandler.DeleteCollaboration)
// 	}

// 	// Donor-specific Routes (e.g., authenticated donors)
// 	donorAuth := router.Group("/api/v1/donor")
// 	{
// 		// Donations
// 		donorAuth.POST("/donate", donationHandler.CreateDonation) // Handles single and split donations
// 		//donorAuth.GET("/donations/:id", donationHandler.GetDonationByID) // Already in public
// 	}

// 	// Admin Routes (can be integrated later with proper auth)
// 	admin := router.Group("/api/v1/admin")
// 	{
// 		admin.POST("/emergency-funds", emergencyFundHandler.CreateEmergencyFund)
// 		admin.PUT("/emergency-funds/:id", emergencyFundHandler.UpdateEmergencyFund)
// 		admin.DELETE("/emergency-funds/:id", emergencyFundHandler.DeleteEmergencyFund)
// 	}

//		// Run the server
//		port := os.Getenv("PORT")
//		if port == "" {
//			port = "8080" // Default port
//		}
//		log.Printf("Server starting on port %s", port)
//		if err := router.Run(":" + port); err != nil {
//			log.Fatalf("Server failed to start: %v", err)
//		}
//	}
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"go-donation-backend/config"
	"go-donation-backend/handlers"
	"go-donation-backend/middleware"
	"go-donation-backend/services"
)

var mongoClient *mongo.Client

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly.")
	}

	// Connect to MongoDB
	var err error
	mongoClient, err = config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ping the DB to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")

	// Initialize services with the MongoDB client
	ngoService := services.NewNGOService(mongoClient)
	userService := services.NewUserService(mongoClient, ngoService)
	donationService := services.NewDonationService(mongoClient)
	expenditureService := services.NewExpenditureService(mongoClient)
	projectUpdateService := services.NewProjectUpdateService(mongoClient)
	certificateService := services.NewCertificateService(mongoClient)
	gamificationService := services.NewGamificationService(mongoClient)
	collaborationService := services.NewCollaborationService(mongoClient)
	emergencyFundService := services.NewEmergencyFundService(mongoClient)

	// Initialize handlers with their respective services
	authHandler := handlers.NewAuthHandler(userService)
	ngoHandler := handlers.NewNGOHandler(ngoService)
	donationHandler := handlers.NewDonationHandler(donationService)
	expenditureHandler := handlers.NewExpenditureHandler(expenditureService)
	projectUpdateHandler := handlers.NewProjectUpdateHandler(projectUpdateService)
	certificateHandler := handlers.NewCertificateHandler(certificateService)
	gamificationHandler := handlers.NewGamificationHandler(gamificationService)
	collaborationHandler := handlers.NewCollaborationHandler(collaborationService)
	emergencyFundHandler := handlers.NewEmergencyFundHandler(emergencyFundService)

	// Setup Gin router
	router := gin.Default()

	// --- CORS Configuration ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- End CORS Configuration ---

	// Auth Routes - Publicly accessible
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Public Routes (e.g., for general viewing by anyone, including unauthenticated users)
	public := router.Group("/api/v1")
	{
		public.GET("/ngos", ngoHandler.GetNGOs)
		public.GET("/ngos/:id", ngoHandler.GetNGOByID)
		public.GET("/ngos/search", ngoHandler.SearchNGOs)
		public.GET("/emergency-funds", emergencyFundHandler.GetEmergencyFunds)
		public.GET("/emergency-funds/:id", emergencyFundHandler.GetEmergencyFundByID)
		public.GET("/collaborations", collaborationHandler.GetCollaborations)
		public.GET("/gamification/leaderboard", gamificationHandler.GetLeaderboard)
		public.GET("/project-updates/ngo/:ngoID", projectUpdateHandler.GetProjectUpdatesByNGO)
		public.GET("/certificates/:id", certificateHandler.GetDonationCertificate)
		public.GET("/certificates/donation/:donationID", certificateHandler.GetCertificateByDonationID)
		public.GET("/donations/ngo/:ngoID", donationHandler.GetDonationsByNGO)
		public.GET("/donations/:id", donationHandler.GetDonationByID)
	}

	// Authenticated Routes - require a valid JWT
	// Create a *gin.RouterGroup for authenticated routes and apply the AuthMiddleware to it.
	// This allows subsequent .Group() calls on 'authenticatedGroup'.
	authenticatedGroup := router.Group("/api/v1")
	authenticatedGroup.Use(middleware.AuthMiddleware())
	{
		// Donor-specific Routes
		donorAuth := authenticatedGroup.Group("/donor") // This is now valid
		{
			donorAuth.POST("/donate", donationHandler.CreateDonation)
			donorAuth.GET("/donations/donor/:donorID", middleware.DonorRequired(), donationHandler.GetDonationsByDonor)
			donorAuth.GET("/gamification/donor/:donorID/achievements", middleware.DonorRequired(), gamificationHandler.GetDonorAchievements)
		}

		// NGO-specific Routes (Requires NGO role and matching NGOID if present in path)
		ngoAuth := authenticatedGroup.Group("/ngo") // This is now valid
		{
			ngoAuth.PUT("/:id", middleware.NGORequired(), ngoHandler.UpdateNGO)
			ngoAuth.DELETE("/:id", middleware.NGORequired(), ngoHandler.DeleteNGO)

			// Expenditures
			ngoAuth.POST("/:ngoID/expenditures", middleware.NGORequired(), expenditureHandler.AddExpenditure)
			ngoAuth.GET("/:ngoID/expenditures", middleware.NGORequired(), expenditureHandler.GetExpendituresByNGO)

			// Project Updates
			ngoAuth.POST("/:ngoID/project-updates", middleware.NGORequired(), projectUpdateHandler.CreateProjectUpdate)
			ngoAuth.PUT("/project-updates/:id", middleware.NGORequired(), projectUpdateHandler.UpdateProjectUpdate)
			ngoAuth.DELETE("/project-updates/:id", middleware.NGORequired(), projectUpdateHandler.DeleteProjectUpdate)

			// Collaborations (NGOs can create/manage their collaborations)
			ngoAuth.POST("/collaborations", middleware.NGORequired(), collaborationHandler.CreateCollaboration)
			ngoAuth.PUT("/collaborations/:id", middleware.NGORequired(), collaborationHandler.UpdateCollaboration)
			ngoAuth.DELETE("/collaborations/:id", middleware.NGORequired(), collaborationHandler.DeleteCollaboration)
			ngoAuth.GET("/collaborations/ngo/:ngoID", middleware.NGORequired(), collaborationHandler.GetCollaborationsByNGO)
		}

		// Admin can manage any NGO (no specific NGOID match required in middleware)
		adminNGOManagement := authenticatedGroup.Group("/admin/ngos").Use(middleware.AdminRequired()) // This is now valid
		{
			adminNGOManagement.POST("/register", ngoHandler.CreateNGO)
			adminNGOManagement.PUT("/:id", ngoHandler.UpdateNGO)
			adminNGOManagement.DELETE("/:id", ngoHandler.DeleteNGO)
		}

		// Admin Routes
		adminAuth := authenticatedGroup.Group("/admin").Use(middleware.AdminRequired()) // This is now valid
		{
			adminAuth.POST("/emergency-funds", emergencyFundHandler.CreateEmergencyFund)
			adminAuth.PUT("/emergency-funds/:id", emergencyFundHandler.UpdateEmergencyFund)
			adminAuth.DELETE("/emergency-funds/:id", emergencyFundHandler.DeleteEmergencyFund)
		}
	}

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
