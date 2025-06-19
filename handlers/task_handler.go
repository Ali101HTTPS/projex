package handlers

import (
	"net/http"
	"project-x/models"
	"project-x/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var createTaskRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		ProjectID   *uint      `json:"project_id"`
		DueDate     *time.Time `json:"due_date"`
		AssignedTo  *uint      `json:"assigned_to"` // Optional: assign to specific user (Admin/Manager only)
	}

	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user info from context
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	// Determine who the task should be assigned to
	var assignedUserID uint
	if createTaskRequest.AssignedTo != nil {
		// Check if current user has permission to assign tasks to others
		if userRole == models.RoleEmployee || userRole == models.RoleHead {
			c.JSON(http.StatusForbidden, gin.H{"error": "Employees and Heads cannot assign tasks to other users"})
			return
		}
		assignedUserID = *createTaskRequest.AssignedTo
	} else {
		// Assign to current user
		assignedUserID = userID.(uint)
	}

	taskService := services.NewTaskService(h.DB)
	task, err := taskService.CreateTask(
		createTaskRequest.Title,
		createTaskRequest.Description,
		assignedUserID,
		createTaskRequest.ProjectID,
		createTaskRequest.DueDate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task": gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"user_id":     task.UserID,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		},
	})
}

// CreateTaskForUser allows Admin/Manager to create tasks for specific users
func (h *TaskHandler) CreateTaskForUser(c *gin.Context) {
	var createTaskRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		UserID      uint       `json:"user_id" binding:"required"`
		ProjectID   *uint      `json:"project_id"`
		DueDate     *time.Time `json:"due_date"`
		Priority    string     `json:"priority"` // high, medium, low
	}

	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user info from context
	currentUserID, _ := c.Get("userID")

	// Validate priority
	if createTaskRequest.Priority != "" {
		validPriorities := []string{"high", "medium", "low"}
		validPriority := false
		for _, p := range validPriorities {
			if p == createTaskRequest.Priority {
				validPriority = true
				break
			}
		}
		if !validPriority {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority. Must be 'high', 'medium', or 'low'"})
			return
		}
	}

	taskService := services.NewTaskService(h.DB)
	task, err := taskService.CreateTaskForUser(
		createTaskRequest.Title,
		createTaskRequest.Description,
		createTaskRequest.UserID,
		currentUserID.(uint),
		createTaskRequest.ProjectID,
		createTaskRequest.DueDate,
		createTaskRequest.Priority,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task assigned successfully",
		"task": gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"user_id":     task.UserID,
			"assigned_by": currentUserID,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		},
	})
}

// CreateCollaborativeTask creates a new collaborative task
func (h *TaskHandler) CreateCollaborativeTask(c *gin.Context) {
	var createTaskRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		ProjectID   *uint      `json:"project_id"`
		DueDate     *time.Time `json:"due_date"`
		AssignedTo  *uint      `json:"assigned_to"` // Optional: assign to specific user (Admin/Manager only)
	}

	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user info from context
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	// Determine who the task should be assigned to
	var assignedUserID uint
	if createTaskRequest.AssignedTo != nil {
		// Check if current user has permission to assign tasks to others
		if userRole == models.RoleEmployee || userRole == models.RoleHead {
			c.JSON(http.StatusForbidden, gin.H{"error": "Employees and Heads cannot assign collaborative tasks to other users"})
			return
		}
		assignedUserID = *createTaskRequest.AssignedTo
	} else {
		// Assign to current user
		assignedUserID = userID.(uint)
	}

	taskService := services.NewTaskService(h.DB)
	task, err := taskService.CreateCollaborativeTask(
		createTaskRequest.Title,
		createTaskRequest.Description,
		assignedUserID,
		createTaskRequest.ProjectID,
		createTaskRequest.DueDate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Collaborative task created successfully",
		"task": gin.H{
			"id":           task.ID,
			"title":        task.Title,
			"description":  task.Description,
			"status":       task.Status,
			"lead_user_id": task.LeadUserID,
			"project_id":   task.ProjectID,
			"assigned_at":  task.AssignedAt,
			"due_date":     task.DueDate,
			"created_at":   task.CreatedAt,
		},
	})
}

// CreateCollaborativeTaskForUser allows Admin/Manager to create collaborative tasks for specific users
func (h *TaskHandler) CreateCollaborativeTaskForUser(c *gin.Context) {
	var createTaskRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		UserID      uint       `json:"user_id" binding:"required"`
		ProjectID   *uint      `json:"project_id"`
		DueDate     *time.Time `json:"due_date"`
		Priority    string     `json:"priority"` // high, medium, low
	}

	if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user info from context
	currentUserID, _ := c.Get("userID")

	// Validate priority
	if createTaskRequest.Priority != "" {
		validPriorities := []string{"high", "medium", "low"}
		validPriority := false
		for _, p := range validPriorities {
			if p == createTaskRequest.Priority {
				validPriority = true
				break
			}
		}
		if !validPriority {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority. Must be 'high', 'medium', or 'low'"})
			return
		}
	}

	taskService := services.NewTaskService(h.DB)
	task, err := taskService.CreateCollaborativeTaskForUser(
		createTaskRequest.Title,
		createTaskRequest.Description,
		createTaskRequest.UserID,
		currentUserID.(uint),
		createTaskRequest.ProjectID,
		createTaskRequest.DueDate,
		createTaskRequest.Priority,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Collaborative task assigned successfully",
		"task": gin.H{
			"id":           task.ID,
			"title":        task.Title,
			"description":  task.Description,
			"status":       task.Status,
			"lead_user_id": task.LeadUserID,
			"assigned_by":  currentUserID,
			"project_id":   task.ProjectID,
			"assigned_at":  task.AssignedAt,
			"due_date":     task.DueDate,
			"created_at":   task.CreatedAt,
		},
	})
}

// GetUserTasks returns all tasks for the current user
func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetUserTasks(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

// GetUserCollaborativeTasks returns all collaborative tasks for the current user
func (h *TaskHandler) GetUserCollaborativeTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetUserCollaborativeTasks(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collaborative tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":           task.ID,
			"title":        task.Title,
			"description":  task.Description,
			"status":       task.Status,
			"lead_user_id": task.LeadUserID,
			"project_id":   task.ProjectID,
			"assigned_at":  task.AssignedAt,
			"due_date":     task.DueDate,
			"created_at":   task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"collaborative_tasks": taskList})
}

// GetProjectTasks returns all tasks for a specific project (Admin or project member)
func (h *TaskHandler) GetProjectTasks(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetProjectTasks(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"user_id":     task.UserID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

// GetProjectCollaborativeTasks returns all collaborative tasks for a specific project
func (h *TaskHandler) GetProjectCollaborativeTasks(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetProjectCollaborativeTasks(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project collaborative tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":             task.ID,
			"title":          task.Title,
			"description":    task.Description,
			"status":         task.Status,
			"lead_user_id":   task.LeadUserID,
			"lead_user_name": task.LeadUser.Username,
			"project_id":     task.ProjectID,
			"assigned_at":    task.AssignedAt,
			"due_date":       task.DueDate,
			"created_at":     task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"collaborative_tasks": taskList})
}

// UpdateTaskStatus updates task status
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updateRequest struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	validStatus := false
	for _, status := range validStatuses {
		if status == updateRequest.Status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	err = taskService.UpdateTaskStatus(uint(taskID), models.TaskStatus(updateRequest.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully"})
}

// UpdateCollaborativeTaskStatus updates collaborative task status
func (h *TaskHandler) UpdateCollaborativeTaskStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updateRequest struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	validStatus := false
	for _, status := range validStatuses {
		if status == updateRequest.Status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	err = taskService.UpdateCollaborativeTaskStatus(uint(taskID), models.TaskStatus(updateRequest.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collaborative task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborative task status updated successfully"})
}

// GetTasksByStatus returns tasks filtered by status for current user
func (h *TaskHandler) GetTasksByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status parameter is required"})
		return
	}

	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetTasksByStatus(userID.(uint), models.TaskStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

// GetCollaborativeTasksByStatus returns collaborative tasks filtered by status for current user
func (h *TaskHandler) GetCollaborativeTasksByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status parameter is required"})
		return
	}

	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetCollaborativeTasksByStatus(userID.(uint), models.TaskStatus(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collaborative tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"collaborative_tasks": taskList})
}

// DeleteTask deletes a task (only by owner)
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	err = taskService.DeleteTask(uint(taskID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// DeleteCollaborativeTask deletes a collaborative task (only by owner)
func (h *TaskHandler) DeleteCollaborativeTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, _ := c.Get("userID")

	taskService := services.NewTaskService(h.DB)
	err = taskService.DeleteCollaborativeTask(uint(taskID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborative task deleted successfully"})
}

// GetTasksByDepartment returns all tasks for users in a specific department (Head/Admin only)
func (h *TaskHandler) GetTasksByDepartment(c *gin.Context) {
	department := c.Param("dept")
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department parameter is required"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetTasksByDepartment(department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"user_id":     task.UserID,
			"user_name":   task.User.Username,
			"project_id":  task.ProjectID,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tasks": taskList})
}

// GetCollaborativeTasksByDepartment returns all collaborative tasks for users in a specific department (Head/Admin only)
func (h *TaskHandler) GetCollaborativeTasksByDepartment(c *gin.Context) {
	department := c.Param("dept")
	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department parameter is required"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	tasks, err := taskService.GetCollaborativeTasksByDepartment(department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department collaborative tasks"})
		return
	}

	var taskList []gin.H
	for _, task := range tasks {
		taskList = append(taskList, gin.H{
			"id":             task.ID,
			"title":          task.Title,
			"description":    task.Description,
			"status":         task.Status,
			"lead_user_id":   task.LeadUserID,
			"lead_user_name": task.LeadUser.Username,
			"project_id":     task.ProjectID,
			"assigned_at":    task.AssignedAt,
			"due_date":       task.DueDate,
			"created_at":     task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"collaborative_tasks": taskList})
}

// BulkUpdateTaskStatus allows Manager/Admin to update multiple task statuses at once
func (h *TaskHandler) BulkUpdateTaskStatus(c *gin.Context) {
	var bulkUpdateRequest struct {
		TaskIDs []uint            `json:"task_ids" binding:"required"`
		Status  models.TaskStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&bulkUpdateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []models.TaskStatus{
		models.TaskStatusPending,
		models.TaskStatusInProgress,
		models.TaskStatusCompleted,
		models.TaskStatusCancelled,
	}
	validStatus := false
	for _, status := range validStatuses {
		if status == bulkUpdateRequest.Status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	updatedCount, err := taskService.BulkUpdateTaskStatus(bulkUpdateRequest.TaskIDs, bulkUpdateRequest.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Tasks updated successfully",
		"updated_count":   updatedCount,
		"total_requested": len(bulkUpdateRequest.TaskIDs),
	})
}

// GetTaskStatistics returns task statistics for managers and admins
func (h *TaskHandler) GetTaskStatistics(c *gin.Context) {
	taskService := services.NewTaskService(h.DB)
	stats, err := taskService.GetTaskStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

// GetProjectReport returns detailed project report with user performance
func (h *TaskHandler) GetProjectReport(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "weekly" // Default to weekly
	}

	if period != "weekly" && period != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Use 'weekly' or 'monthly'"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	report, err := taskService.GetProjectReport(uint(projectID), period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetUserReport returns detailed user report with project performance
func (h *TaskHandler) GetUserReport(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "weekly" // Default to weekly
	}

	if period != "weekly" && period != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Use 'weekly' or 'monthly'"})
		return
	}

	taskService := services.NewTaskService(h.DB)
	report, err := taskService.GetUserReport(uint(userID), period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}
