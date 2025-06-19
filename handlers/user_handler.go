package handlers

import (
	"net/http"
	"project-x/models"
	"project-x/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// CreateUser creates a new user (Admin only)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserRequest struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		Role       string `json:"role" binding:"required"`
		Department string `json:"department" binding:"required"`
	}

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.CreateUser(
		createUserRequest.Username,
		createUserRequest.Password,
		createUserRequest.Role,
		createUserRequest.Department,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
			"created_at": user.CreatedAt,
		},
	})
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// ListUsers returns all users (Admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
	userService := services.NewUserService(h.DB)
	users, err := userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userList})
}

// GetUsersByRole returns users filtered by role
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	userService := services.NewUserService(h.DB)
	users, err := userService.GetUsersByRole(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userList})
}

// GetUsersByDepartment returns users filtered by department
func (h *UserHandler) GetUsersByDepartment(c *gin.Context) {
	department := c.Param("department")
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department parameter is required"})
		return
	}

	userService := services.NewUserService(h.DB)
	users, err := userService.GetUsersByDepartment(department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userList})
}

// UpdateUserRole updates a user's role (Admin only)
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateRequest struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.UpdateUserRole(uint(userID), updateRequest.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
		},
	})
}

// UpdateUserDepartment updates a user's department (Admin only)
func (h *UserHandler) UpdateUserDepartment(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateRequest struct {
		Department string `json:"department" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService(h.DB)
	user, err := userService.UpdateUserDepartment(uint(userID), updateRequest.Department)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User department updated successfully",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
		},
	})
}

// UpdateUserPassword updates a user's password (Admin or self)
func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is updating their own password or is admin
	currentUserID, _ := c.Get("userID")
	currentUserRole, _ := c.Get("userRole")

	if currentUserID.(uint) != uint(userID) && currentUserRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own password"})
		return
	}

	var updateRequest struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService(h.DB)
	err = userService.UpdateUserPassword(uint(userID), updateRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// GetUserStats returns statistics about a user
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := services.NewUserService(h.DB)
	stats, err := userService.GetUserStats(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// DeleteUser deletes a user (Admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := services.NewUserService(h.DB)
	err = userService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
