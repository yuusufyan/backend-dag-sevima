package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	model "backend-sevima/internal/models"
	"backend-sevima/internal/routing"

	"github.com/joho/godotenv"

	"github.com/yuusufyan/go-common/pkg/database"
	"github.com/yuusufyan/go-common/pkg/logger"
	commonfiber "github.com/yuusufyan/go-common/pkg/middleware/fiber"
	"github.com/yuusufyan/go-common/pkg/utils"
)

// @title SEVIMA Backend API
// @version 1.0
// @description This is the API documentation for SEVIMA Backend.
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load .env if it exists
	_ = godotenv.Load()

	// Initialize Logger
	isProd := os.Getenv("APP_ENV") == "production"
	appLogger := logger.New(isProd)

	// Initialize Database Connection
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if dbPort == 0 {
		dbPort = 5432
	}

	dbConfig := &database.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     dbPort,
	}

	// For local testing without valid DB credentials, we don't want to panic immediately.
	// But in a real app, we usually fatal here. Let's log an error if it fails, but continue so the server starts.
	db, err := database.Connect(dbConfig, appLogger, isProd)
	if err != nil {
		appLogger.WithError(err).Warn("Failed to connect to database. DAG Engine will not work until DB is configured.")
	} else {
		appLogger.Info("Successfully connected to database")

		// AutoMigrate Database Models
		err = db.AutoMigrate(
			&model.Tenant{},
			&model.User{},
			&model.DagTemplate{},
			&model.DagExecution{},
			&model.TaskInstance{},
		)
		if err != nil {
			appLogger.WithError(err).Error("Failed to auto migrate database models")
		} else {
			appLogger.Info("Database models migrated successfully")
		}
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(commonfiber.Telemetry())
	app.Use(commonfiber.Recover(appLogger))
	app.Use(commonfiber.Logger(appLogger))

	// Use go-common middleware for identity (extracts X-User-ID, etc)
	app.Use(commonfiber.Identity())

	// Initialize Redis (if not used, it will gracefully return nil)
	rdb := database.InitRedis(&database.RedisConfig{
		Host: "", // Kosongkan agar tidak memaksa connect jika belum ada
	}, appLogger)

	// Register Standard Health Check on root app (unauthenticated)
	utils.RegisterHealthCheck(app, db, rdb)

	// API Routes
	routing.SetupRoutes(app, db)

	// Start server
	go func() {
		if err := app.Listen(":3000"); err != nil {
			appLogger.WithError(err).Panic("Fiber listen error")
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		appLogger.WithError(err).Fatal("Server forced to shutdown")
	}

	appLogger.Info("Server exiting")
}
