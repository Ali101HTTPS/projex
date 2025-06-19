package routes

import (
	"project-x/handlers"
	"project-x/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTaskRoutes(r *gin.Engine, db *gorm.DB) {
	taskHandler := handlers.NewTaskHandler(db)

	taskGroup := r.Group("/api/tasks")
	taskGroup.Use(middleware.AuthMiddleware(db))
	{
		// Basic task endpoints - Head/Manager/Admin only (exclude Employee)
		taskGroup.POST("", middleware.RequireHeadOrHigher(), taskHandler.CreateTask)
		taskGroup.GET("", taskHandler.GetUserTasks) // All users can view their own tasks
		taskGroup.GET("/status/:status", taskHandler.GetTasksByStatus)
		taskGroup.PATCH("/:id/status", taskHandler.UpdateTaskStatus)                       // All users can update their own task status
		taskGroup.DELETE("/:id", middleware.RequireHeadOrHigher(), taskHandler.DeleteTask) // Head+ can delete tasks

		// Collaborative task endpoints - Head/Manager/Admin only (exclude Employee)
		taskGroup.POST("/collaborative", middleware.RequireHeadOrHigher(), taskHandler.CreateCollaborativeTask)
		taskGroup.GET("/collaborative", taskHandler.GetUserCollaborativeTasks)
		taskGroup.GET("/collaborative/status/:status", taskHandler.GetCollaborativeTasksByStatus)
		taskGroup.PATCH("/collaborative/:id/status", taskHandler.UpdateCollaborativeTaskStatus)
		taskGroup.DELETE("/collaborative/:id", middleware.RequireHeadOrHigher(), taskHandler.DeleteCollaborativeTask) // Head+ can delete collaborative tasks

		// Project task endpoints - project members only
		taskGroup.GET("/project/:projectId", taskHandler.GetProjectTasks)
		taskGroup.GET("/project/:projectId/collaborative", taskHandler.GetProjectCollaborativeTasks)

		// Admin/Manager task assignment endpoints
		taskGroup.POST("/assign", middleware.RequireManagerOrHigher(), taskHandler.CreateTaskForUser)
		taskGroup.POST("/collaborative/assign", middleware.RequireManagerOrHigher(), taskHandler.CreateCollaborativeTaskForUser)

		// Department-based endpoints - Manager/Admin only (since Manager > Head)
		taskGroup.GET("/department/:dept", middleware.RequireManagerOrHigher(), taskHandler.GetTasksByDepartment)
		taskGroup.GET("/department/:dept/collaborative", middleware.RequireManagerOrHigher(), taskHandler.GetCollaborativeTasksByDepartment)

		// Management endpoints - Manager+ can access
		taskGroup.GET("/all", middleware.RequireManagerOrHigher(), taskHandler.GetUserTasks)                  // Manager+ can see all tasks
		taskGroup.PATCH("/:id/assign", middleware.RequireManagerOrHigher(), taskHandler.UpdateTaskStatus)     // Manager+ can assign tasks
		taskGroup.POST("/bulk-update", middleware.RequireManagerOrHigher(), taskHandler.BulkUpdateTaskStatus) // Manager+ can bulk update
		taskGroup.GET("/statistics", middleware.RequireManagerOrHigher(), taskHandler.GetTaskStatistics)      // Manager+ can see statistics

		// Report endpoints - Manager+ can access
		taskGroup.GET("/reports/project/:projectId", middleware.RequireManagerOrHigher(), taskHandler.GetProjectReport) // Project report
		taskGroup.GET("/reports/user/:userId", middleware.RequireManagerOrHigher(), taskHandler.GetUserReport)          // User report

		// Admin-only endpoints
		taskGroup.DELETE("/:id/force", middleware.RequireAdmin(), taskHandler.DeleteTask)         // Admin only can force delete
		taskGroup.PATCH("/:id/reassign", middleware.RequireAdmin(), taskHandler.UpdateTaskStatus) // Admin only can reassign
	}
}
