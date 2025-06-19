package routes

import (
	"project-x/handlers"
	"project-x/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCollaborativeTaskRoutes(r *gin.Engine, db *gorm.DB) {
	collaborativeTaskHandler := handlers.NewCollaborativeTaskHandler(db)

	collaborativeTaskGroup := r.Group("/api/collaborative-tasks")
	collaborativeTaskGroup.Use(middleware.AuthMiddleware(db))
	{
		// Create collaborative task (Head/Manager/Admin only)
		collaborativeTaskGroup.POST("", middleware.RequireHeadOrHigher(), collaborativeTaskHandler.CreateCollaborativeTask)

		// Get user's collaborative tasks (all participants)
		collaborativeTaskGroup.GET("", collaborativeTaskHandler.GetUserCollaborativeTasks)

		// Get collaborative task details (participants only)
		collaborativeTaskGroup.GET("/:id", collaborativeTaskHandler.GetCollaborativeTaskDetails)

		// Get collaborative task statistics (participants only)
		collaborativeTaskGroup.GET("/:id/statistics", collaborativeTaskHandler.GetCollaborativeTaskStatistics)

		// Add participant to collaborative task (Lead/Manager/Admin only)
		collaborativeTaskGroup.POST("/:id/participants", middleware.RequireHeadOrHigher(), collaborativeTaskHandler.AddParticipant)

		// Remove participant from collaborative task (Lead/Manager/Admin only)
		collaborativeTaskGroup.DELETE("/:id/participants/:userId", middleware.RequireHeadOrHigher(), collaborativeTaskHandler.RemoveParticipant)

		// Update task progress (participants only)
		collaborativeTaskGroup.PATCH("/:id/progress", collaborativeTaskHandler.UpdateTaskProgress)
	}
}
