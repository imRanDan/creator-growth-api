package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/api"
	"github.com/imRanDan/creator-growth-api/internal/database"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	//Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file foind")
	}

	//Connect to the database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	//Run the migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	//start background job
	services.StartTokenRefreshJob()

	//Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"message":  "Creator Growth API is running",
			"database": "connected",
		})
	})

	//rate limiter setup (5 attemps per minute peer IP)
	rate, _ := limiter.NewRateFromFormatted("5-M")
	store := memory.NewStore()
	limitedMiddleware := limiter.New(store, rate).Middleware()

	//Pub routes
	r.POST("/api/auth/register", api.Register)
	r.POST("/api/auth/login", limitedMiddleware, api.Login) // rate limit on login
	r.GET("/auth/instagram/callback", api.InstagramCallback)

	//Protected Routes (requires auth/JWT)
	protected := r.Group("/api")
	protected.Use(api.AuthMiddleware())
	{
		protected.GET("/user/me", api.GetCurrentUser)
		protected.GET("/instagram/connect", api.ConnectInstagram)
		protected.GET("/growth/stats", api.GetGrowthStats)
	}

	//Start Server
	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
