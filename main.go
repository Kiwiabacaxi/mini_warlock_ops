package main

import (
	"log"
	"os"

	"github-analytics-backend/github"
	"github-analytics-backend/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Initialize GitHub client
	githubClient := github.NewClient(token)

	// Initialize handlers
	handler := handlers.NewHandler(githubClient)

	// Setup Gin router
	r := gin.Default()

	// Enable CORS for frontend debugging
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Repository endpoints
		api.GET("/repositories/:owner", handler.GetRepositories)

		// Metrics endpoints
		api.GET("/metrics/commits/:owner/:repo", handler.GetCommitMetrics)
		api.GET("/metrics/prs/:owner/:repo", handler.GetPRMetrics)
		api.GET("/metrics/contributions/:owner/:repo", handler.GetContributionMetrics)

		// Chart endpoints
		api.GET("/charts/commits-leaderboard/:owner/:repo", handler.GetCommitsLeaderboard)
		api.GET("/charts/commits-timeline/:owner/:repo", handler.GetCommitsTimeline)
		api.GET("/charts/prs-timeline/:owner/:repo", handler.GetPRsTimeline)
		api.GET("/charts/prs-leaderboard/:owner/:repo", handler.GetPRsLeaderboard)
		api.GET("/charts/contributions-leaderboard/:owner/:repo", handler.GetContributionsLeaderboard)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "GitHub Analytics Backend is running",
		})
	})

	// Serve static files (frontend)
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Printf("Frontend available at: http://localhost:%s", port)
	log.Printf("API available at: http://localhost:%s/api", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
