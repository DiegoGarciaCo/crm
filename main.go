package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/DiegoGarciaCo/CRM/internal/handlers"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/joho/godotenv"
	"github.com/keighl/postmark"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	// ------------------------------------------------
	// Get configuration from environment variables
	// ------------------------------------------------

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}
	devMode := os.Getenv("DEV")
	dev := devMode == "true"

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}
	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("S3REGION is not set")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3BUCKET is not set")
	}
	postmarkServerToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	if postmarkServerToken == "" {
		log.Fatal("POSTMARK_SERVER_TOKEN is not set")
	}
	betterAuthSecret := os.Getenv("BETTER_AUTH_SECRET")
	if betterAuthSecret == "" {
		log.Fatal("BETTER_AUTH_SECRET is not set")
	}
	serverURL := os.Getenv("SERVER_BASE_URL")
	if serverURL == "" {
		log.Fatal("SERVER_BASE_URL is not set")
	}
	fromEmail := os.Getenv("FROM_EMAIL_ADDRESS")
	if fromEmail == "" {
		log.Fatal("FROM_EMAIL_ADDRESS is not set")
	}

	// -----------------------------------------------
	// Initialize Logger
	// -----------------------------------------------

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	// ------------------------------------------------
	// Initialize database connection
	// ------------------------------------------------

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Check if the database is reachable
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	dbQueries := database.New(db)
	if err != nil {
		log.Fatalf("Error creating database queries: %v", err)
	}

	// ------------------------------------------------
	// Initialize S3 client
	// ------------------------------------------------

	awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s3Region))
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(awsConfig)

	// ------------------------------------------------
	// Initialize Postmark client
	// ------------------------------------------------

	postmarkClient := postmark.Client{
		ServerToken: postmarkServerToken,
		BaseURL:     "https://api.postmarkapp.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	EmailSecret := []byte(JWTSecret)

	// ------------------------------------------------
	// Initialize config, server and cors
	// ------------------------------------------------
	cfg := handlers.New(port, JWTSecret, dbQueries, db, dev, logger, s3Client, s3Bucket, s3Region, &postmarkClient, EmailSecret, betterAuthSecret, serverURL, fromEmail)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://access.soldbyghost.com", "https://app.soldbyghost.com", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-KEY"},
		AllowCredentials: true,
	})

	// Create a new HTTP server mux
	mux := http.NewServeMux()

	// ------------------------------------------------
	// Define routes and handlers
	// ------------------------------------------------

	// Dashboard Routes
	mux.HandleFunc("GET /api/dashboard/new-contacts", cfg.GetNewContactsCount)
	mux.HandleFunc("GET /api/dashboard/appointments", cfg.GetAppointmentsCount)
	mux.HandleFunc("GET /api/dashboard/tasks-today", cfg.GetTasksDueTodayCount)
	mux.HandleFunc("GET /api/dashboard/5-newest-contacts", cfg.Get5NewestContacts)
	mux.HandleFunc("GET /api/dashboard/5-upcoming-appointments", cfg.Get5UpcomingAppointments)
	mux.HandleFunc("GET /api/dashboard/contacts-count", cfg.GetContactsCount)
	mux.HandleFunc("GET /api/dashboard/contacts-by-source", cfg.ContactCountBySource)

	// Contact Routes
	mux.HandleFunc("POST /api/contacts", cfg.CreateContact)
	mux.HandleFunc("POST /api/contacts/import", cfg.ImportContacts)
	mux.HandleFunc("GET /api/contacts/contact/{contactID}", cfg.GetContactByID)
	mux.HandleFunc("GET /api/contacts", cfg.GetAllContacts)
	mux.HandleFunc("GET /api/contacts/search", cfg.SearchContacts)
	mux.HandleFunc("GET /api/contacts/smart-list/{smartListID}", cfg.GetContactsBySmartList)

	// Notes Routes
	mux.HandleFunc("POST /api/notes", cfg.CreateNote)
	mux.HandleFunc("GET /api/notes/{contactID}", cfg.GetNotesByContactID)

	// Contact Logs Routes
	mux.HandleFunc("POST /api/contact-logs", cfg.LogContact)
	mux.HandleFunc("GET /api/contact-logs/{contactID}", cfg.GetContactLogsByContactID)

	// Tasks Routes
	mux.HandleFunc("POST /api/tasks", cfg.CreateTask)
	mux.HandleFunc("GET /api/tasks/contact/{contactID}", cfg.GetTasksByContactID)
	mux.HandleFunc("GET /api/tasks/assigned", cfg.GetTaskByAssignedToID)
	mux.HandleFunc("GET /api/tasks/{taskID}", cfg.GetTaskByID)
	mux.HandleFunc("DELETE /api/tasks/{taskID}", cfg.DeleteTask)
	mux.HandleFunc("PUT /api/tasks/{taskID}", cfg.UpdateTask)
	mux.HandleFunc("GET /api/tasks/late", cfg.GetOverdueTasks)
	mux.HandleFunc("PUT /api/tasks/status/{taskID}", cfg.UpdateTaskStatus)
	mux.HandleFunc("GET /api/tasks/today", cfg.GetTasksDueToday)

	// Appointments Routes
	mux.HandleFunc("POST /api/appointments", cfg.CreateAppointment)
	mux.HandleFunc("GET /api/appointments/{AppointmentID}", cfg.GetAppointmentByID)
	mux.HandleFunc("PUT /api/appointments/{AppointmentID}", cfg.UpdateAppointment)
	mux.HandleFunc("DELETE /api/appointments/{AppointmentID}", cfg.DeleteAppointment)
	mux.HandleFunc("GET /api/appointments/contact/{ContactID}", cfg.ListAppointmentsByContactID)
	mux.HandleFunc("GET /api/appointments/upcoming", cfg.ListUpcomingAppointments)
	mux.HandleFunc("GET /api/appointments/today", cfg.ListAppointmentsToday)
	mux.HandleFunc("GET /api/appointments", cfg.ListAppointments)

	// Deals Routes
	mux.HandleFunc("POST /api/deals", cfg.CreateDeal)
	mux.HandleFunc("GET /api/deals/{dealID}", cfg.GetDealByID)
	mux.HandleFunc("PUT /api/deals/{dealID}", cfg.UpdateDeal)
	mux.HandleFunc("DELETE /api/deals/{dealID}", cfg.DeleteDeal)
	mux.HandleFunc("GET /api/deals", cfg.ListDeals)
	mux.HandleFunc("GET /api/deals/contacts/{contactID}", cfg.ListDealsByContactID)
	mux.HandleFunc("GET /api/deals/stage/{stageID}", cfg.ListDealsByStageID)

	// Goals Routes
	mux.HandleFunc("POST /api/goals", cfg.SetGoal)
	mux.HandleFunc("GET /api/goals", cfg.GetGoalByUserAndYear)
	mux.HandleFunc("PUT /api/goals/{GoalID}", cfg.UpdateGoal)

	// Smart Lists Routes
	mux.HandleFunc("GET /api/smart-lists", cfg.GetAllSmartLists)
	mux.HandleFunc("POST /api/smart-lists", cfg.CreateSmartList)
	mux.HandleFunc("PUT /api/smart-lists/{smartListID}/filter", cfg.SetSmartListFilterCriteria)
	mux.HandleFunc("PUT /api/smart-lists/{smartListID}/name", cfg.UpdateSmartList)

	// Stages Routes
	mux.HandleFunc("POST /api/stages", cfg.CreateStage)
	mux.HandleFunc("GET /api/stages", cfg.GetStages)
	mux.HandleFunc("GET /api/stages/client-type", cfg.GetStagesByClientType)

	// Tags Routes
	mux.HandleFunc("POST /api/tags", cfg.CreateTag)
	mux.HandleFunc("GET /api/tags", cfg.GetAllTags)
	mux.HandleFunc("DELETE /api/tags/{tagID}", cfg.DeleteTag)
	mux.HandleFunc("POST /api/tags/{tagID}/contact/{contactID}", cfg.AssignTagToContact)
	mux.HandleFunc("DELETE /api/tags/{tagID}/contact/{contactID}", cfg.RemoveTagFromContact)

	// Webhooks Routes
	mux.HandleFunc("POST /webhooks/landing-page-form", cfg.CollectLandingPageForm)

	// Email Routes
	mux.HandleFunc("GET /api/verify", cfg.VerifyEmail)
	mux.HandleFunc("POST /api/resend-verification", cfg.ResendVerificationEmail)

	// S3 Routes
	mux.HandleFunc("PUT /api/upload-profile-picture", cfg.UploadProfilePicture)

	// ------------------------------------------------------------
	// Wrap mux with CORS handler, middleware and start server
	// ------------------------------------------------------------
	handler := cfg.AuthMiddleware()(mux)
	handler = cfg.LoggerMiddleware(handler)
	handler = corsHandler.Handler(handler)

	// Configure server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	log.Print("Listening on port " + cfg.Port + " ...")
	if err := srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}
}
