package handlers

import (
	"net/http"
	"project-x/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CollaborativeTaskHandler struct {
	DB *gorm.DB
}

func NewCollaborativeTaskHandler(db *gorm.DB) *CollaborativeTaskHandler {
	return &CollaborativeTaskHandler{DB: db}
}

// CreateCollaborativeTask creates a new collaborative task (Head/Manager/Admin only)
func (h *CollaborativeTaskHandler) CreateCollaborativeTask(c *gin.Context) {
	var createRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		LeadUserID  uint       `json:"lead_user_id" binding:"required"`
		ProjectID   *uint      `json:"project_id"`
		DueDate     *time.Time `json:"due_date"`
		Priority    string     `json:"priority" binding:"required"`
		Complexity  string     `json:"complexity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	task, err := collaborativeTaskService.CreateCollaborativeTask(
		createRequest.Title,
		createRequest.Description,
		createRequest.LeadUserID,
		createRequest.ProjectID,
		createRequest.DueDate,
		createRequest.Priority,
		createRequest.Complexity,
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
			"priority":     task.Priority,
			"complexity":   task.Complexity,
			"progress":     task.Progress,
			"assigned_at":  task.AssignedAt,
			"due_date":     task.DueDate,
			"created_at":   task.CreatedAt,
		},
	})
}

// AddParticipant adds a user to a collaborative task
func (h *CollaborativeTaskHandler) AddParticipant(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var addRequest struct {
		UserID       uint   `json:"user_id" binding:"required"`
		Role         string `json:"role" binding:"required"`
		Contribution string `json:"contribution"`
	}

	if err := c.ShouldBindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	err = collaborativeTaskService.AddParticipant(uint(taskID), addRequest.UserID, addRequest.Role, addRequest.Contribution)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant added successfully"})
}

// RemoveParticipant removes a user from a collaborative task
func (h *CollaborativeTaskHandler) RemoveParticipant(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	err = collaborativeTaskService.RemoveParticipant(uint(taskID), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}

// UpdateTaskProgress updates the progress of a collaborative task
func (h *CollaborativeTaskHandler) UpdateTaskProgress(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updateRequest struct {
		Progress int `json:"progress" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	err = collaborativeTaskService.UpdateTaskProgress(uint(taskID), updateRequest.Progress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task progress updated successfully"})
}

// GetCollaborativeTaskDetails returns detailed information about a collaborative task
func (h *CollaborativeTaskHandler) GetCollaborativeTaskDetails(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	task, err := collaborativeTaskService.GetCollaborativeTaskWithDetails(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collaborative task not found"})
		return
	}

	// Build participants list
	var participants []gin.H
	for _, participant := range task.Participants {
		participants = append(participants, gin.H{
			"id":           participant.ID,
			"user_id":      participant.UserID,
			"username":     participant.User.Username,
			"role":         participant.Role,
			"status":       participant.Status,
			"contribution": participant.Contribution,
			"assigned_at":  participant.AssignedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"task": gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"lead_user": gin.H{
				"id":       task.LeadUser.ID,
				"username": task.LeadUser.Username,
			},
			"project_id":   task.ProjectID,
			"priority":     task.Priority,
			"complexity":   task.Complexity,
			"progress":     task.Progress,
			"assigned_at":  task.AssignedAt,
			"due_date":     task.DueDate,
			"created_at":   task.CreatedAt,
			"participants": participants,
		},
	})
}

// GetCollaborativeTaskStatistics returns statistics for a collaborative task
func (h *CollaborativeTaskHandler) GetCollaborativeTaskStatistics(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	stats, err := collaborativeTaskService.GetCollaborativeTaskStatistics(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}

// GetUserCollaborativeTasks returns all collaborative tasks where user is a participant
func (h *CollaborativeTaskHandler) GetUserCollaborativeTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	collaborativeTaskService := services.NewCollaborativeTaskService(h.DB)
	tasks, err := collaborativeTaskService.GetUserCollaborativeTasks(userID.(uint))
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
			"lead_user": gin.H{
				"id":       task.LeadUser.ID,
				"username": task.LeadUser.Username,
			},
			"project_id":  task.ProjectID,
			"priority":    task.Priority,
			"complexity":  task.Complexity,
			"progress":    task.Progress,
			"assigned_at": task.AssignedAt,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"collaborative_tasks": taskList})
}
