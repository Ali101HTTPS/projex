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

type ProjectHandler struct {
	DB *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{DB: db}
}

// CreateProject creates a new project (Manager/Admin only)
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var createProjectRequest struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description" binding:"required"`
		StartDate   time.Time  `json:"start_date" binding:"required"`
		EndDate     *time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&createProjectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user ID from context
	userID, _ := c.Get("userID")

	projectService := services.NewProjectService(h.DB)
	project, err := projectService.CreateProject(
		createProjectRequest.Title,
		createProjectRequest.Description,
		userID.(uint),
		createProjectRequest.StartDate,
		createProjectRequest.EndDate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": gin.H{
			"id":          project.ID,
			"title":       project.Title,
			"description": project.Description,
			"status":      project.Status,
			"created_by":  project.CreatedBy,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"created_at":  project.CreatedAt,
		},
	})
}

// GetUserProjects returns all projects the current user is part of
func (h *ProjectHandler) GetUserProjects(c *gin.Context) {
	userID, _ := c.Get("userID")

	projectService := services.NewProjectService(h.DB)
	projects, err := projectService.GetUserProjects(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	var projectList []gin.H
	for _, project := range projects {
		projectList = append(projectList, gin.H{
			"id":          project.ID,
			"title":       project.Title,
			"description": project.Description,
			"status":      project.Status,
			"created_by":  project.CreatedBy,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"created_at":  project.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectList})
}

// GetProjectDetails returns detailed information about a specific project
func (h *ProjectHandler) GetProjectDetails(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectService := services.NewProjectService(h.DB)
	project, err := projectService.GetProjectWithDetails(uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": gin.H{
			"id":          project.ID,
			"title":       project.Title,
			"description": project.Description,
			"status":      project.Status,
			"created_by":  project.CreatedBy,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"created_at":  project.CreatedAt,
			"creator": gin.H{
				"id":       project.Creator.ID,
				"username": project.Creator.Username,
				"role":     project.Creator.Role,
			},
			"member_count": len(project.Users),
			"task_count":   len(project.Tasks) + len(project.CollaborativeTasks),
		},
	})
}

// GetProjectMembers returns all members of a project
func (h *ProjectHandler) GetProjectMembers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectService := services.NewProjectService(h.DB)
	members, err := projectService.GetProjectMembers(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project members"})
		return
	}

	var memberList []gin.H
	for _, member := range members {
		memberList = append(memberList, gin.H{
			"id":         member.ID,
			"username":   member.Username,
			"role":       member.Role,
			"department": member.Department,
		})
	}

	c.JSON(http.StatusOK, gin.H{"members": memberList})
}

// AddUserToProject adds a user to a project (Manager/Admin only)
func (h *ProjectHandler) AddUserToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var addUserRequest struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&addUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectService := services.NewProjectService(h.DB)
	err = projectService.AddUserToProject(addUserRequest.UserID, uint(projectID), addUserRequest.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to project successfully"})
}

// RemoveUserFromProject removes a user from a project (Manager/Admin only)
func (h *ProjectHandler) RemoveUserFromProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	projectService := services.NewProjectService(h.DB)
	err = projectService.RemoveUserFromProject(uint(userID), uint(projectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from project successfully"})
}

// UpdateProjectStatus updates project status (Manager/Admin only)
func (h *ProjectHandler) UpdateProjectStatus(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
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
	validStatuses := []string{"active", "paused", "completed", "cancelled"}
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

	projectService := services.NewProjectService(h.DB)
	err = projectService.UpdateProjectStatus(uint(projectID), models.ProjectStatus(updateRequest.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project status updated successfully"})
}

// DeleteProject deletes a project (Admin only)
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectService := services.NewProjectService(h.DB)
	err = projectService.DeleteProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetProjectStatistics returns statistics for a project
func (h *ProjectHandler) GetProjectStatistics(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get project details
	projectService := services.NewProjectService(h.DB)
	project, err := projectService.GetProjectWithDetails(uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Count tasks by status
	var pendingTasks, inProgressTasks, completedTasks, cancelledTasks int
	for _, task := range project.Tasks {
		switch task.Status {
		case models.TaskStatusPending:
			pendingTasks++
		case models.TaskStatusInProgress:
			inProgressTasks++
		case models.TaskStatusCompleted:
			completedTasks++
		case models.TaskStatusCancelled:
			cancelledTasks++
		}
	}

	for _, task := range project.CollaborativeTasks {
		switch task.Status {
		case models.TaskStatusPending:
			pendingTasks++
		case models.TaskStatusInProgress:
			inProgressTasks++
		case models.TaskStatusCompleted:
			completedTasks++
		case models.TaskStatusCancelled:
			cancelledTasks++
		}
	}

	totalTasks := len(project.Tasks) + len(project.CollaborativeTasks)
	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	stats := gin.H{
		"project_id":        project.ID,
		"project_title":     project.Title,
		"status":            project.Status,
		"member_count":      len(project.Users),
		"total_tasks":       totalTasks,
		"pending_tasks":     pendingTasks,
		"in_progress_tasks": inProgressTasks,
		"completed_tasks":   completedTasks,
		"cancelled_tasks":   cancelledTasks,
		"completion_rate":   completionRate,
		"start_date":        project.StartDate,
		"end_date":          project.EndDate,
	}

	c.JSON(http.StatusOK, gin.H{"statistics": stats})
}
