package routes

import (
	"project-x/handlers"
	"project-x/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProjectRoutes(router *gin.Engine, db *gorm.DB) {
	projectHandler := handlers.NewProjectHandler(db)

	// Project routes group
	projects := router.Group("/api/projects")
	{
		// Create project (Manager/Admin only)
		projects.POST("/", middleware.AuthMiddleware(db), middleware.RequireManagerOrHigher(), projectHandler.CreateProject)

		// Get user's projects (all authenticated users)
		projects.GET("/my-projects", middleware.AuthMiddleware(db), projectHandler.GetUserProjects)

		// Get project details (project members only)
		projects.GET("/:id", middleware.AuthMiddleware(db), projectHandler.GetProjectDetails)

		// Get project members (project members only)
		projects.GET("/:id/members", middleware.AuthMiddleware(db), projectHandler.GetProjectMembers)

		// Add user to project (Manager/Admin only)
		projects.POST("/:id/members", middleware.AuthMiddleware(db), middleware.RequireManagerOrHigher(), projectHandler.AddUserToProject)

		// Remove user from project (Manager/Admin only)
		projects.DELETE("/:id/members/:userId", middleware.AuthMiddleware(db), middleware.RequireManagerOrHigher(), projectHandler.RemoveUserFromProject)

		// Update project status (Manager/Admin only)
		projects.PATCH("/:id/status", middleware.AuthMiddleware(db), middleware.RequireManagerOrHigher(), projectHandler.UpdateProjectStatus)

		// Get project statistics (project members only)
		projects.GET("/:id/statistics", middleware.AuthMiddleware(db), projectHandler.GetProjectStatistics)

		// Delete project (Admin only)
		projects.DELETE("/:id", middleware.AuthMiddleware(db), middleware.RequireAdmin(), projectHandler.DeleteProject)
	}
}
