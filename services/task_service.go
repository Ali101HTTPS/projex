package services

import (
	"errors"
	"project-x/models"
	"time"

	"gorm.io/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{DB: db}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(title, description string, userID uint, projectID *uint, dueDate *time.Time) (*models.Task, error) {
	// Verify user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// If projectID is provided, verify project exists and user is member
	if projectID != nil {
		var userProject models.UserProject
		if err := s.DB.Where("user_id = ? AND project_id = ?", userID, *projectID).First(&userProject).Error; err != nil {
			return nil, errors.New("user is not a member of this project")
		}
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		UserID:      userID,
		ProjectID:   projectID,
		AssignedAt:  time.Now(),
		DueDate:     dueDate,
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// CreateTaskForUser allows Admin/Manager to create tasks for specific users
func (s *TaskService) CreateTaskForUser(title, description string, userID, assignedBy uint, projectID *uint, dueDate *time.Time, priority string) (*models.Task, error) {
	// Verify target user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("target user not found")
	}

	// Verify assigner exists and has appropriate role
	var assigner models.User
	if err := s.DB.First(&assigner, assignedBy).Error; err != nil {
		return nil, errors.New("assigner not found")
	}

	// Check if assigner has permission (Admin or Manager only)
	if assigner.Role == models.RoleEmployee || assigner.Role == models.RoleHead {
		return nil, errors.New("employees and heads cannot assign tasks to other users")
	}

	// If projectID is provided, verify project exists and user is member
	if projectID != nil {
		var userProject models.UserProject
		if err := s.DB.Where("user_id = ? AND project_id = ?", userID, *projectID).First(&userProject).Error; err != nil {
			return nil, errors.New("target user is not a member of this project")
		}
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		UserID:      userID,
		ProjectID:   projectID,
		AssignedAt:  time.Now(),
		DueDate:     dueDate,
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// CreateCollaborativeTask creates a new collaborative task
func (s *TaskService) CreateCollaborativeTask(title, description string, userID uint, projectID *uint, dueDate *time.Time) (*models.CollaborativeTask, error) {
	// Verify user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// If projectID is provided, verify project exists and user is member
	if projectID != nil {
		var userProject models.UserProject
		if err := s.DB.Where("user_id = ? AND project_id = ?", userID, *projectID).First(&userProject).Error; err != nil {
			return nil, errors.New("user is not a member of this project")
		}
	}

	task := &models.CollaborativeTask{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		LeadUserID:  userID,
		ProjectID:   projectID,
		AssignedAt:  time.Now(),
		DueDate:     dueDate,
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// CreateCollaborativeTaskForUser allows Admin/Manager to create collaborative tasks for specific users
func (s *TaskService) CreateCollaborativeTaskForUser(title, description string, userID, assignedBy uint, projectID *uint, dueDate *time.Time, priority string) (*models.CollaborativeTask, error) {
	// Verify target user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("target user not found")
	}

	// Verify assigner exists and has appropriate role
	var assigner models.User
	if err := s.DB.First(&assigner, assignedBy).Error; err != nil {
		return nil, errors.New("assigner not found")
	}

	// Check if assigner has permission (Admin or Manager only)
	if assigner.Role == models.RoleEmployee || assigner.Role == models.RoleHead {
		return nil, errors.New("employees and heads cannot assign collaborative tasks to other users")
	}

	// If projectID is provided, verify project exists and user is member
	if projectID != nil {
		var userProject models.UserProject
		if err := s.DB.Where("user_id = ? AND project_id = ?", userID, *projectID).First(&userProject).Error; err != nil {
			return nil, errors.New("target user is not a member of this project")
		}
	}

	task := &models.CollaborativeTask{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		LeadUserID:  userID,
		ProjectID:   projectID,
		AssignedAt:  time.Now(),
		DueDate:     dueDate,
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// GetUserTasks returns all tasks for a user
func (s *TaskService) GetUserTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := s.DB.Where("user_id = ?", userID).
		Preload("Project").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetUserCollaborativeTasks returns all collaborative tasks for a user
func (s *TaskService) GetUserCollaborativeTasks(userID uint) ([]models.CollaborativeTask, error) {
	var tasks []models.CollaborativeTask
	err := s.DB.Where("user_id = ?", userID).
		Preload("Project").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetProjectTasks returns all tasks for a project
func (s *TaskService) GetProjectTasks(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := s.DB.Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetProjectCollaborativeTasks returns all collaborative tasks for a project
func (s *TaskService) GetProjectCollaborativeTasks(projectID uint) ([]models.CollaborativeTask, error) {
	var tasks []models.CollaborativeTask
	err := s.DB.Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// UpdateTaskStatus updates task status
func (s *TaskService) UpdateTaskStatus(taskID uint, status models.TaskStatus) error {
	return s.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("status", status).Error
}

// UpdateCollaborativeTaskStatus updates collaborative task status
func (s *TaskService) UpdateCollaborativeTaskStatus(taskID uint, status models.TaskStatus) error {
	return s.DB.Model(&models.CollaborativeTask{}).Where("id = ?", taskID).Update("status", status).Error
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(taskID, userID uint) error {
	var task models.Task
	if err := s.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return errors.New("task not found or access denied")
	}

	return s.DB.Delete(&task).Error
}

// DeleteCollaborativeTask deletes a collaborative task
func (s *TaskService) DeleteCollaborativeTask(taskID, userID uint) error {
	var task models.CollaborativeTask
	if err := s.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return errors.New("task not found or access denied")
	}

	return s.DB.Delete(&task).Error
}

// GetTasksByStatus returns tasks filtered by status
func (s *TaskService) GetTasksByStatus(userID uint, status models.TaskStatus) ([]models.Task, error) {
	var tasks []models.Task
	err := s.DB.Where("user_id = ? AND status = ?", userID, status).
		Preload("Project").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetCollaborativeTasksByStatus returns collaborative tasks filtered by status
func (s *TaskService) GetCollaborativeTasksByStatus(userID uint, status models.TaskStatus) ([]models.CollaborativeTask, error) {
	var tasks []models.CollaborativeTask
	err := s.DB.Where("user_id = ? AND status = ?", userID, status).
		Preload("Project").
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetTasksByDepartment returns all tasks for users in a specific department
func (s *TaskService) GetTasksByDepartment(department string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.DB.Joins("JOIN users ON tasks.user_id = users.id").
		Where("users.department = ?", department).
		Preload("User").
		Preload("Project").
		Order("tasks.created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// GetCollaborativeTasksByDepartment returns all collaborative tasks for users in a specific department
func (s *TaskService) GetCollaborativeTasksByDepartment(department string) ([]models.CollaborativeTask, error) {
	var tasks []models.CollaborativeTask
	err := s.DB.Joins("JOIN users ON collaborative_tasks.user_id = users.id").
		Where("users.department = ?", department).
		Preload("User").
		Preload("Project").
		Order("collaborative_tasks.created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// BulkUpdateTaskStatus updates multiple task statuses at once
func (s *TaskService) BulkUpdateTaskStatus(taskIDs []uint, status models.TaskStatus) (int64, error) {
	result := s.DB.Model(&models.Task{}).Where("id IN ?", taskIDs).Update("status", status)
	return result.RowsAffected, result.Error
}

// GetTaskStatistics returns overall task statistics
func (s *TaskService) GetTaskStatistics() (map[string]interface{}, error) {
	var totalTasks, pendingTasks, inProgressTasks, completedTasks, cancelledTasks int64
	var totalCollaborativeTasks, pendingCollaborativeTasks, inProgressCollaborativeTasks, completedCollaborativeTasks, cancelledCollaborativeTasks int64

	// Count regular tasks by status
	s.DB.Model(&models.Task{}).Count(&totalTasks)
	s.DB.Model(&models.Task{}).Where("status = ?", models.TaskStatusPending).Count(&pendingTasks)
	s.DB.Model(&models.Task{}).Where("status = ?", models.TaskStatusInProgress).Count(&inProgressTasks)
	s.DB.Model(&models.Task{}).Where("status = ?", models.TaskStatusCompleted).Count(&completedTasks)
	s.DB.Model(&models.Task{}).Where("status = ?", models.TaskStatusCancelled).Count(&cancelledTasks)

	// Count collaborative tasks by status
	s.DB.Model(&models.CollaborativeTask{}).Count(&totalCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("status = ?", models.TaskStatusPending).Count(&pendingCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("status = ?", models.TaskStatusInProgress).Count(&inProgressCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("status = ?", models.TaskStatusCompleted).Count(&completedCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("status = ?", models.TaskStatusCancelled).Count(&cancelledCollaborativeTasks)

	// Calculate completion rates
	totalAllTasks := totalTasks + totalCollaborativeTasks
	totalCompletedTasks := completedTasks + completedCollaborativeTasks

	completionRate := 0.0
	if totalAllTasks > 0 {
		completionRate = float64(totalCompletedTasks) / float64(totalAllTasks) * 100
	}

	stats := map[string]interface{}{
		"regular_tasks": map[string]interface{}{
			"total":       totalTasks,
			"pending":     pendingTasks,
			"in_progress": inProgressTasks,
			"completed":   completedTasks,
			"cancelled":   cancelledTasks,
		},
		"collaborative_tasks": map[string]interface{}{
			"total":       totalCollaborativeTasks,
			"pending":     pendingCollaborativeTasks,
			"in_progress": inProgressCollaborativeTasks,
			"completed":   completedCollaborativeTasks,
			"cancelled":   cancelledCollaborativeTasks,
		},
		"overall": map[string]interface{}{
			"total_tasks":     totalAllTasks,
			"completed_tasks": totalCompletedTasks,
			"completion_rate": completionRate,
		},
	}

	return stats, nil
}

// GetProjectReport returns detailed statistics for a specific project
func (s *TaskService) GetProjectReport(projectID uint, period string) (map[string]interface{}, error) {
	var startDate, endDate time.Time
	now := time.Now()

	// Calculate date range based on period
	switch period {
	case "weekly":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "monthly":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	default:
		return nil, errors.New("invalid period. Use 'weekly' or 'monthly'")
	}

	// Get project details
	var project models.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return nil, errors.New("project not found")
	}

	// Regular tasks statistics
	var totalTasks, pendingTasks, inProgressTasks, completedTasks, cancelledTasks int64
	var tasksInPeriod int64

	s.DB.Model(&models.Task{}).Where("project_id = ?", projectID).Count(&totalTasks)
	s.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusPending).Count(&pendingTasks)
	s.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusInProgress).Count(&inProgressTasks)
	s.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusCompleted).Count(&completedTasks)
	s.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusCancelled).Count(&cancelledTasks)
	s.DB.Model(&models.Task{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startDate, endDate).Count(&tasksInPeriod)

	// Collaborative tasks statistics
	var totalCollaborativeTasks, pendingCollaborativeTasks, inProgressCollaborativeTasks, completedCollaborativeTasks, cancelledCollaborativeTasks int64
	var collaborativeTasksInPeriod int64

	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ?", projectID).Count(&totalCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusPending).Count(&pendingCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusInProgress).Count(&inProgressCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusCompleted).Count(&completedCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND status = ?", projectID, models.TaskStatusCancelled).Count(&cancelledCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND created_at BETWEEN ? AND ?", projectID, startDate, endDate).Count(&collaborativeTasksInPeriod)

	// User performance in this project
	var userStats []map[string]interface{}
	var users []models.User
	s.DB.Joins("JOIN user_projects ON users.id = user_projects.user_id").
		Where("user_projects.project_id = ?", projectID).
		Find(&users)

	for _, user := range users {
		var userCompletedTasks, userTotalTasks int64
		var userCompletedCollaborativeTasks, userTotalCollaborativeTasks int64

		// Count user's regular tasks in this project
		s.DB.Model(&models.Task{}).Where("project_id = ? AND user_id = ?", projectID, user.ID).Count(&userTotalTasks)
		s.DB.Model(&models.Task{}).Where("project_id = ? AND user_id = ? AND status = ?", projectID, user.ID, models.TaskStatusCompleted).Count(&userCompletedTasks)

		// Count user's collaborative tasks in this project (as lead)
		s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND lead_user_id = ?", projectID, user.ID).Count(&userTotalCollaborativeTasks)
		s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND lead_user_id = ? AND status = ?", projectID, user.ID, models.TaskStatusCompleted).Count(&userCompletedCollaborativeTasks)

		userCompletionRate := 0.0
		totalUserTasks := userTotalTasks + userTotalCollaborativeTasks
		totalUserCompleted := userCompletedTasks + userCompletedCollaborativeTasks
		if totalUserTasks > 0 {
			userCompletionRate = float64(totalUserCompleted) / float64(totalUserTasks) * 100
		}

		userStats = append(userStats, map[string]interface{}{
			"user_id":         user.ID,
			"username":        user.Username,
			"role":            user.Role,
			"department":      user.Department,
			"total_tasks":     totalUserTasks,
			"completed_tasks": totalUserCompleted,
			"completion_rate": userCompletionRate,
			"regular_tasks": map[string]interface{}{
				"total":     userTotalTasks,
				"completed": userCompletedTasks,
			},
			"collaborative_tasks": map[string]interface{}{
				"total":     userTotalCollaborativeTasks,
				"completed": userCompletedCollaborativeTasks,
			},
		})
	}

	// Calculate overall project completion rate
	totalProjectTasks := totalTasks + totalCollaborativeTasks
	totalProjectCompleted := completedTasks + completedCollaborativeTasks
	projectCompletionRate := 0.0
	if totalProjectTasks > 0 {
		projectCompletionRate = float64(totalProjectCompleted) / float64(totalProjectTasks) * 100
	}

	report := map[string]interface{}{
		"project": map[string]interface{}{
			"id":          project.ID,
			"title":       project.Title,
			"description": project.Description,
			"status":      project.Status,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
		},
		"period": map[string]interface{}{
			"type":       period,
			"start_date": startDate,
			"end_date":   endDate,
		},
		"statistics": map[string]interface{}{
			"regular_tasks": map[string]interface{}{
				"total":             totalTasks,
				"pending":           pendingTasks,
				"in_progress":       inProgressTasks,
				"completed":         completedTasks,
				"cancelled":         cancelledTasks,
				"created_in_period": tasksInPeriod,
			},
			"collaborative_tasks": map[string]interface{}{
				"total":             totalCollaborativeTasks,
				"pending":           pendingCollaborativeTasks,
				"in_progress":       inProgressCollaborativeTasks,
				"completed":         completedCollaborativeTasks,
				"cancelled":         cancelledCollaborativeTasks,
				"created_in_period": collaborativeTasksInPeriod,
			},
			"overall": map[string]interface{}{
				"total_tasks":     totalProjectTasks,
				"completed_tasks": totalProjectCompleted,
				"completion_rate": projectCompletionRate,
			},
		},
		"user_performance": userStats,
	}

	return report, nil
}

// GetUserReport returns detailed statistics for a specific user
func (s *TaskService) GetUserReport(userID uint, period string) (map[string]interface{}, error) {
	var startDate, endDate time.Time
	now := time.Now()

	// Calculate date range based on period
	switch period {
	case "weekly":
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case "monthly":
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	default:
		return nil, errors.New("invalid period. Use 'weekly' or 'monthly'")
	}

	// Get user details
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Regular tasks statistics
	var totalTasks, pendingTasks, inProgressTasks, completedTasks, cancelledTasks int64
	var tasksInPeriod, completedInPeriod int64

	s.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&totalTasks)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusPending).Count(&pendingTasks)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusInProgress).Count(&inProgressTasks)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusCompleted).Count(&completedTasks)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusCancelled).Count(&cancelledTasks)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).Count(&tasksInPeriod)
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ? AND updated_at BETWEEN ? AND ?", userID, models.TaskStatusCompleted, startDate, endDate).Count(&completedInPeriod)

	// Collaborative tasks statistics (as lead)
	var totalCollaborativeTasks, pendingCollaborativeTasks, inProgressCollaborativeTasks, completedCollaborativeTasks, cancelledCollaborativeTasks int64
	var collaborativeTasksInPeriod, collaborativeCompletedInPeriod int64

	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ?", userID).Count(&totalCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND status = ?", userID, models.TaskStatusPending).Count(&pendingCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND status = ?", userID, models.TaskStatusInProgress).Count(&inProgressCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND status = ?", userID, models.TaskStatusCompleted).Count(&completedCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND status = ?", userID, models.TaskStatusCancelled).Count(&cancelledCollaborativeTasks)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).Count(&collaborativeTasksInPeriod)
	s.DB.Model(&models.CollaborativeTask{}).Where("lead_user_id = ? AND status = ? AND updated_at BETWEEN ? AND ?", userID, models.TaskStatusCompleted, startDate, endDate).Count(&collaborativeCompletedInPeriod)

	// Project performance breakdown
	var projectStats []map[string]interface{}
	var projects []models.Project
	s.DB.Joins("JOIN user_projects ON projects.id = user_projects.project_id").
		Where("user_projects.user_id = ?", userID).
		Find(&projects)

	for _, project := range projects {
		var projectTasks, projectCompleted int64
		var projectCollaborativeTasks, projectCollaborativeCompleted int64

		// Count user's tasks in this project
		s.DB.Model(&models.Task{}).Where("project_id = ? AND user_id = ?", project.ID, userID).Count(&projectTasks)
		s.DB.Model(&models.Task{}).Where("project_id = ? AND user_id = ? AND status = ?", project.ID, userID, models.TaskStatusCompleted).Count(&projectCompleted)

		// Count user's collaborative tasks in this project (as lead)
		s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND lead_user_id = ?", project.ID, userID).Count(&projectCollaborativeTasks)
		s.DB.Model(&models.CollaborativeTask{}).Where("project_id = ? AND lead_user_id = ? AND status = ?", project.ID, userID, models.TaskStatusCompleted).Count(&projectCollaborativeCompleted)

		projectCompletionRate := 0.0
		totalProjectUserTasks := projectTasks + projectCollaborativeTasks
		totalProjectUserCompleted := projectCompleted + projectCollaborativeCompleted
		if totalProjectUserTasks > 0 {
			projectCompletionRate = float64(totalProjectUserCompleted) / float64(totalProjectUserTasks) * 100
		}

		projectStats = append(projectStats, map[string]interface{}{
			"project_id":      project.ID,
			"project_title":   project.Title,
			"project_status":  project.Status,
			"total_tasks":     totalProjectUserTasks,
			"completed_tasks": totalProjectUserCompleted,
			"completion_rate": projectCompletionRate,
			"regular_tasks": map[string]interface{}{
				"total":     projectTasks,
				"completed": projectCompleted,
			},
			"collaborative_tasks": map[string]interface{}{
				"total":     projectCollaborativeTasks,
				"completed": projectCollaborativeCompleted,
			},
		})
	}

	// Calculate overall user completion rate
	totalUserTasks := totalTasks + totalCollaborativeTasks
	totalUserCompleted := completedTasks + completedCollaborativeTasks
	userCompletionRate := 0.0
	if totalUserTasks > 0 {
		userCompletionRate = float64(totalUserCompleted) / float64(totalUserTasks) * 100
	}

	// Calculate period performance
	totalInPeriod := tasksInPeriod + collaborativeTasksInPeriod
	totalCompletedInPeriod := completedInPeriod + collaborativeCompletedInPeriod
	periodCompletionRate := 0.0
	if totalInPeriod > 0 {
		periodCompletionRate = float64(totalCompletedInPeriod) / float64(totalInPeriod) * 100
	}

	report := map[string]interface{}{
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"department": user.Department,
		},
		"period": map[string]interface{}{
			"type":       period,
			"start_date": startDate,
			"end_date":   endDate,
		},
		"statistics": map[string]interface{}{
			"regular_tasks": map[string]interface{}{
				"total":               totalTasks,
				"pending":             pendingTasks,
				"in_progress":         inProgressTasks,
				"completed":           completedTasks,
				"cancelled":           cancelledTasks,
				"created_in_period":   tasksInPeriod,
				"completed_in_period": completedInPeriod,
			},
			"collaborative_tasks": map[string]interface{}{
				"total":               totalCollaborativeTasks,
				"pending":             pendingCollaborativeTasks,
				"in_progress":         inProgressCollaborativeTasks,
				"completed":           completedCollaborativeTasks,
				"cancelled":           cancelledCollaborativeTasks,
				"created_in_period":   collaborativeTasksInPeriod,
				"completed_in_period": collaborativeCompletedInPeriod,
			},
			"overall": map[string]interface{}{
				"total_tasks":     totalUserTasks,
				"completed_tasks": totalUserCompleted,
				"completion_rate": userCompletionRate,
			},
			"period_performance": map[string]interface{}{
				"total_tasks":     totalInPeriod,
				"completed_tasks": totalCompletedInPeriod,
				"completion_rate": periodCompletionRate,
			},
		},
		"project_performance": projectStats,
	}

	return report, nil
}
