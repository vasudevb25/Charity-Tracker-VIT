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
// 	"go-donation-backend/middleware"
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
// 	organizationService := services.NewOrganizationService(mongoClient) // <--- Updated service
// 	userService := services.NewUserService(mongoClient, organizationService)
// 	donationService := services.NewDonationService(mongoClient)
// 	expenditureService := services.NewExpenditureService(mongoClient)
// 	projectUpdateService := services.NewProjectUpdateService(mongoClient)
// 	certificateService := services.NewCertificateService(mongoClient)
// 	gamificationService := services.NewGamificationService(mongoClient)
// 	collaborationService := services.NewCollaborationService(mongoClient)
// 	reportService := services.NewReportService(mongoClient) // <--- New Report Service

// 	// Initialize handlers with their respective services
// 	authHandler := handlers.NewAuthHandler(userService)
// 	organizationHandler := handlers.NewOrganizationHandler(organizationService) // <--- Updated handler
// 	donationHandler := handlers.NewDonationHandler(donationService)
// 	expenditureHandler := handlers.NewExpenditureHandler(expenditureService)
// 	projectUpdateHandler := handlers.NewProjectUpdateHandler(projectUpdateService)
// 	certificateHandler := handlers.NewCertificateHandler(certificateService)
// 	gamificationHandler := handlers.NewGamificationHandler(gamificationService)
// 	collaborationHandler := handlers.NewCollaborationHandler(collaborationService)
// 	reportHandler := handlers.NewReportHandler(reportService, donationService) // <--- New Report Handler

// 	// Setup Gin router
// 	router := gin.Default()

// 	// --- CORS Configuration ---
// 	router.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))
// 	// --- End CORS Configuration ---

// 	// Auth Routes - Publicly accessible
// 	authRoutes := router.Group("/api/v1/auth")
// 	{
// 		authRoutes.POST("/register", authHandler.Register)
// 		authRoutes.POST("/login", authHandler.Login)
// 	}

// 	// Public Routes (e.g., for general viewing by anyone, including unauthenticated users)
// 	public := router.Group("/api/v1")
// 	{
// 		// Organizations (viewing public organization profiles)
// 		public.GET("/organizations", organizationHandler.GetOrganizations)           // <--- Updated route
// 		public.GET("/organizations/:id", organizationHandler.GetOrganizationByID)    // <--- Updated route
// 		public.GET("/organizations/search", organizationHandler.SearchOrganizations) // <--- Updated route

// 		// Collaborations (viewing public collaborations)
// 		public.GET("/collaborations", collaborationHandler.GetCollaborations)

// 		// Leaderboard (public viewing)
// 		public.GET("/gamification/leaderboard", gamificationHandler.GetLeaderboard)

// 		// Project Updates (public viewing)
// 		public.GET("/organizations/:orgID/project-updates", projectUpdateHandler.GetProjectUpdatesByOrganization) // <--- Updated route

// 		// Donation certificates (public viewing, if ID is known)
// 		public.GET("/certificates/:id", certificateHandler.GetDonationCertificate)
// 		public.GET("/certificates/donation/:donationID", certificateHandler.GetCertificateByDonationID)

// 		// Donations by Organization (public viewing, to show organization's received donations)
// 		public.GET("/organizations/:orgID/donations", donationHandler.GetDonationsByOrganization) // <--- Updated route
// 		public.GET("/donations/:id", donationHandler.GetDonationByID)                             // Get single donation, if ID is known
// 	}

// 	// Authenticated Routes - require a valid JWT
// 	authenticatedGroup := router.Group("/api/v1")
// 	authenticatedGroup.Use(middleware.AuthMiddleware())
// 	{
// 		// Donor-specific Routes
// 		donorAuth := authenticatedGroup.Group("/donor")
// 		{
// 			donorAuth.POST("/donate", donationHandler.CreateDonation)                                                            // Handles single and split donations
// 			donorAuth.GET("/:donorID/donations", middleware.DonorRequired(), donationHandler.GetDonationsByDonor)                // <--- Updated route
// 			donorAuth.GET("/:donorID/achievements", middleware.DonorRequired(), gamificationHandler.GetDonorAchievements)        // <--- Updated route
// 			donorAuth.GET("/:donorID/transaction-history", middleware.DonorRequired(), reportHandler.GetDonorTransactionHistory) // <--- New Route
// 		}

// 		// Organization-specific Routes (Requires Organization role and matching OrganizationID if present in path)
// 		organizationAuth := authenticatedGroup.Group("/organization") // <--- Updated group name
// 		{
// 			organizationAuth.PUT("/:id", middleware.OrganizationRequired(), organizationHandler.UpdateOrganization)    // <--- Updated handler/middleware
// 			organizationAuth.DELETE("/:id", middleware.OrganizationRequired(), organizationHandler.DeleteOrganization) // <--- Updated handler/middleware

// 			// Expenditures
// 			organizationAuth.POST("/:orgID/expenditures", middleware.OrganizationRequired(), expenditureHandler.AddExpenditure)               // <--- Updated handler/middleware
// 			organizationAuth.GET("/:orgID/expenditures", middleware.OrganizationRequired(), expenditureHandler.GetExpendituresByOrganization) // <--- Updated handler/middleware

// 			// Project Updates
// 			organizationAuth.POST("/:orgID/project-updates", middleware.OrganizationRequired(), projectUpdateHandler.CreateProjectUpdate) // <--- Updated handler/middleware
// 			organizationAuth.PUT("/project-updates/:id", middleware.OrganizationRequired(), projectUpdateHandler.UpdateProjectUpdate)     // <--- Updated handler/middleware
// 			organizationAuth.DELETE("/project-updates/:id", middleware.OrganizationRequired(), projectUpdateHandler.DeleteProjectUpdate)  // <--- Updated handler/middleware

// 			// Collaborations (Organizations can create/manage their collaborations)
// 			organizationAuth.POST("/collaborations", middleware.OrganizationRequired(), collaborationHandler.CreateCollaboration)                   // <--- Updated handler/middleware
// 			organizationAuth.PUT("/collaborations/:id", middleware.OrganizationRequired(), collaborationHandler.UpdateCollaboration)                // <--- Updated handler/middleware
// 			organizationAuth.DELETE("/collaborations/:id", middleware.OrganizationRequired(), collaborationHandler.DeleteCollaboration)             // <--- Updated handler/middleware
// 			organizationAuth.GET("/:orgID/collaborations", middleware.OrganizationRequired(), collaborationHandler.GetCollaborationsByOrganization) // <--- Updated handler/middleware

// 			// Organization Audit Report
// 			organizationAuth.GET("/:orgID/audit-report", middleware.OrganizationRequired(), reportHandler.GetOrganizationAuditReport) // <--- New Route
// 		}

// 		// Admin can manage any Organization (no specific OrganizationID match required in middleware)
// 		adminOrganizationManagement := authenticatedGroup.Group("/admin/organizations").Use(middleware.AdminRequired()) // <--- Updated group name
// 		{
// 			adminOrganizationManagement.POST("", organizationHandler.CreateOrganization) // Admin can create new Organizations
// 			adminOrganizationManagement.PUT("/:id", organizationHandler.UpdateOrganization)
// 			adminOrganizationManagement.DELETE("/:id", organizationHandler.DeleteOrganization)
// 		}

// 		// Admin Routes (for system-level management)
// 		// adminAuth := authenticatedGroup.Group("/admin").Use(middleware.AdminRequired())
// 		// {
// 		// 	// No emergency fund related routes here anymore
// 		// 	// Add other admin-specific management routes as needed, e.g., user management, etc.
// 		// }
// 	}

// 	// Run the server
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080" // Default port
// 	}
// 	log.Printf("Server starting on port %s", port)
// 	if err := router.Run(":" + port); err != nil {
// 		log.Fatalf("Server failed to start: %v", err)
// 	}
// }

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
	organizationService := services.NewOrganizationService(mongoClient)
	userService := services.NewUserService(mongoClient, organizationService)
	donationService := services.NewDonationService(mongoClient)
	expenditureService := services.NewExpenditureService(mongoClient)
	projectUpdateService := services.NewProjectUpdateService(mongoClient)
	certificateService := services.NewCertificateService(mongoClient)
	gamificationService := services.NewGamificationService(mongoClient)
	collaborationService := services.NewCollaborationService(mongoClient)
	reportService := services.NewReportService(mongoClient)

	// Initialize handlers with their respective services
	authHandler := handlers.NewAuthHandler(userService)
	organizationHandler := handlers.NewOrganizationHandler(organizationService)
	donationHandler := handlers.NewDonationHandler(donationService)
	expenditureHandler := handlers.NewExpenditureHandler(expenditureService)
	projectUpdateHandler := handlers.NewProjectUpdateHandler(projectUpdateService)
	certificateHandler := handlers.NewCertificateHandler(certificateService)
	gamificationHandler := handlers.NewGamificationHandler(gamificationService)
	collaborationHandler := handlers.NewCollaborationHandler(collaborationService)
	reportHandler := handlers.NewReportHandler(reportService, donationService)

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
		// Organizations (viewing public organization profiles)
		public.GET("/organizations", organizationHandler.GetOrganizations)
		public.GET("/organizations/:id", organizationHandler.GetOrganizationByID)
		public.GET("/organizations/search", organizationHandler.SearchOrganizations)

		// Project Updates for a specific organization
		public.GET("/organizations/:id/project-updates", projectUpdateHandler.GetProjectUpdatesByOrganization) // <--- Corrected to :id

		// Donations by Organization
		public.GET("/organizations/:id/donations", donationHandler.GetDonationsByOrganization) // <--- Corrected to :id
		public.GET("/donations/:id", donationHandler.GetDonationByID)

		// Collaborations (viewing public collaborations)
		public.GET("/collaborations", collaborationHandler.GetCollaborations)

		// Leaderboard (public viewing)
		public.GET("/gamification/leaderboard", gamificationHandler.GetLeaderboard)

		// Donation certificates (public viewing, if ID is known)
		public.GET("/certificates/:id", certificateHandler.GetDonationCertificate)
		public.GET("/certificates/donation/:donationID", certificateHandler.GetCertificateByDonationID)

	}

	// Authenticated Routes - require a valid JWT
	authenticatedGroup := router.Group("/api/v1")
	authenticatedGroup.Use(middleware.AuthMiddleware())
	{
		// Donor-specific Routes
		donorAuth := authenticatedGroup.Group("/donor")
		{
			// Note: The middleware.DonorRequired() might need to check c.Param("id") instead of c.Param("donorID")
			// if you are using :id in the path. Please adjust your middleware/handlers accordingly.
			donorAuth.POST("/donate", donationHandler.CreateDonation)
			donorAuth.GET("/:id/donations", middleware.DonorRequired(), donationHandler.GetDonationsByDonor)                // <--- Corrected to :id
			donorAuth.GET("/:id/achievements", middleware.DonorRequired(), gamificationHandler.GetDonorAchievements)        // <--- Corrected to :id
			donorAuth.GET("/:id/transaction-history", middleware.DonorRequired(), reportHandler.GetDonorTransactionHistory) // <--- Corrected to :id
		}

		// Organization-specific Routes (Requires Organization role and matching OrganizationID if present in path)
		organizationAuth := authenticatedGroup.Group("/organization")
		{
			// All routes here should expect the organization ID as ':id'
			organizationAuth.PUT("/:id", middleware.OrganizationRequired(), organizationHandler.UpdateOrganization)
			organizationAuth.DELETE("/:id", middleware.OrganizationRequired(), organizationHandler.DeleteOrganization)

			// Expenditures
			organizationAuth.POST("/:id/expenditures", middleware.OrganizationRequired(), expenditureHandler.AddExpenditure)               // <--- Corrected to :id
			organizationAuth.GET("/:id/expenditures", middleware.OrganizationRequired(), expenditureHandler.GetExpendituresByOrganization) // <--- Corrected to :id

			// Project Updates
			organizationAuth.POST("/:id/project-updates", middleware.OrganizationRequired(), projectUpdateHandler.CreateProjectUpdate) // <--- Corrected to :id
			organizationAuth.PUT("/project-updates/:id", middleware.OrganizationRequired(), projectUpdateHandler.UpdateProjectUpdate)
			organizationAuth.DELETE("/project-updates/:id", middleware.OrganizationRequired(), projectUpdateHandler.DeleteProjectUpdate)

			// Collaborations (Organizations can create/manage their collaborations)
			organizationAuth.POST("/collaborations", middleware.OrganizationRequired(), collaborationHandler.CreateCollaboration)
			organizationAuth.PUT("/collaborations/:id", middleware.OrganizationRequired(), collaborationHandler.UpdateCollaboration)
			organizationAuth.DELETE("/collaborations/:id", middleware.OrganizationRequired(), collaborationHandler.DeleteCollaboration)
			organizationAuth.GET("/:id/collaborations", middleware.OrganizationRequired(), collaborationHandler.GetCollaborationsByOrganization) // <--- Corrected to :id

			// Organization Audit Report
			organizationAuth.GET("/:id/audit-report", middleware.OrganizationRequired(), reportHandler.GetOrganizationAuditReport) // <--- Corrected to :id
		}

		// Admin can manage any Organization (no specific OrganizationID match required in middleware)
		adminOrganizationManagement := authenticatedGroup.Group("/admin/organizations").Use(middleware.AdminRequired())
		{
			adminOrganizationManagement.POST("", organizationHandler.CreateOrganization)
			adminOrganizationManagement.PUT("/:id", organizationHandler.UpdateOrganization)
			adminOrganizationManagement.DELETE("/:id", organizationHandler.DeleteOrganization)
		}

		// Admin Routes (for system-level management) - Example for other admin routes
		// adminAuth := authenticatedGroup.Group("/admin").Use(middleware.AdminRequired())
		// {
		// 	// Add other admin-specific management routes as needed, e.g., user management, etc.
		// }
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
