package routes

import (
	"project-x/handlers"
	"project-x/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)

	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(db))
	{
		// Admin only routes
		userGroup.POST("", middleware.RequireAdmin(), userHandler.CreateUser)
		userGroup.GET("", middleware.RequireAdmin(), userHandler.ListUsers)
		userGroup.GET("/role/:role", middleware.RequireAdmin(), userHandler.GetUsersByRole)
		userGroup.GET("/department/:department", middleware.RequireAdmin(), userHandler.GetUsersByDepartment)
		userGroup.PATCH("/:id/role", middleware.RequireAdmin(), userHandler.UpdateUserRole)
		userGroup.PATCH("/:id/department", middleware.RequireAdmin(), userHandler.UpdateUserDepartment)
		userGroup.DELETE("/:id", middleware.RequireAdmin(), userHandler.DeleteUser)

		// Routes accessible by admin or the user themselves
		userGroup.GET("/:id", middleware.RequireSelfOrAdmin(), userHandler.GetUser)
		userGroup.GET("/:id/stats", middleware.RequireSelfOrAdmin(), userHandler.GetUserStats)
		userGroup.PATCH("/:id/password", middleware.RequireSelfOrAdmin(), userHandler.UpdateUserPassword)
	}
}
