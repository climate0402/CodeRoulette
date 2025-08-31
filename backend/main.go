package main

import (
	"log"
	"os"

	"coderoulette/internal/config"
	"coderoulette/internal/database"
	"coderoulette/internal/handlers"
	"coderoulette/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient, err := database.InitializeRedis(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize services
	matchService := services.NewMatchService(redisClient)
	matchService.SetDB(db)
	problemService := services.NewProblemService(db)
	judgeService := services.NewJudgeService()
	judgeService.SetDB(db)
	reportService := services.NewReportService(db)
	skillCardService := services.NewSkillCardService(redisClient)

	// Initialize handlers
	handlers := handlers.NewHandlers(
		matchService,
		problemService,
		judgeService,
		reportService,
		skillCardService,
	)

	// Setup routes
	router := gin.Default()
	handlers.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
