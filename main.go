package main

import (
	"fmt"
	"log"
	"os"
	"project-x/config"
	"project-x/models"
	"project-x/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Setup database connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate database tables
	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.CollaborativeTask{}, &models.CollaborativeTaskParticipant{}, &models.Project{}, &models.UserProject{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database tables migrated successfully")

	// Initialize routes
	setupRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server starting on port %s", port)
	r.Run(":" + port)
}

func setupDatabase() (*gorm.DB, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	routes.SetupAuthRoutes(r, db)
	routes.SetupUserRoutes(r, db)
	routes.SetupTaskRoutes(r, db)
	routes.SetupProjectRoutes(r, db)
	routes.SetupCollaborativeTaskRoutes(r, db)
}
