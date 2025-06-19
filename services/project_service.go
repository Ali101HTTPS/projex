package services

import (
	"errors"
	"project-x/models"
	"time"

	"gorm.io/gorm"
)

type ProjectService struct {
	DB *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{DB: db}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(title, description string, createdBy uint, startDate time.Time, endDate *time.Time) (*models.Project, error) {
	project := &models.Project{
		Title:       title,
		Description: description,
		Status:      models.ProjectStatusActive,
		CreatedBy:   createdBy,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.DB.Create(project).Error; err != nil {
		return nil, err
	}

	// Add the creator as a member with manager role
	userProject := &models.UserProject{
		ProjectID: project.ID,
		UserID:    createdBy,
		Role:      "manager",
		JoinedAt:  time.Now(),
	}

	if err := s.DB.Create(userProject).Error; err != nil {
		return nil, err
	}

	return project, nil
}

// GetUserProjects returns all projects a user is part of
func (s *ProjectService) GetUserProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := s.DB.Joins("JOIN user_projects ON projects.id = user_projects.project_id").
		Where("user_projects.user_id = ?", userID).
		Preload("Creator").
		Find(&projects).Error
	return projects, err
}

// GetProjectWithDetails returns a project with all related data
func (s *ProjectService) GetProjectWithDetails(projectID uint) (*models.Project, error) {
	var project models.Project
	err := s.DB.Preload("Creator").
		Preload("Users").
		Preload("Tasks").
		Preload("CollaborativeTasks").
		First(&project, projectID).Error
	return &project, err
}

// GetProjectMembers returns all members of a project
func (s *ProjectService) GetProjectMembers(projectID uint) ([]models.User, error) {
	var users []models.User
	err := s.DB.Joins("JOIN user_projects ON users.id = user_projects.user_id").
		Where("user_projects.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}

// AddUserToProject adds a user to a project
func (s *ProjectService) AddUserToProject(userID, projectID uint, role string) error {
	// Check if project exists
	var project models.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return errors.New("project not found")
	}

	// Check if user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if user is already in the project
	var existingUserProject models.UserProject
	if err := s.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&existingUserProject).Error; err == nil {
		return errors.New("user is already a member of this project")
	}

	// Validate role
	validRoles := []string{"manager", "head", "employee"}
	validRole := false
	for _, r := range validRoles {
		if r == role {
			validRole = true
			break
		}
	}
	if !validRole {
		return errors.New("invalid role")
	}

	userProject := &models.UserProject{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  time.Now(),
	}

	return s.DB.Create(userProject).Error
}

// RemoveUserFromProject removes a user from a project
func (s *ProjectService) RemoveUserFromProject(userID, projectID uint) error {
	// Check if user is the project creator
	var project models.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return errors.New("project not found")
	}

	if project.CreatedBy == userID {
		return errors.New("cannot remove project creator")
	}

	// Remove user from project
	result := s.DB.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.UserProject{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user is not a member of this project")
	}

	return nil
}

// UpdateProjectStatus updates the status of a project
func (s *ProjectService) UpdateProjectStatus(projectID uint, status models.ProjectStatus) error {
	return s.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("status", status).Error
}

// DeleteProject deletes a project and all related data
func (s *ProjectService) DeleteProject(projectID uint) error {
	// Start a transaction
	tx := s.DB.Begin()

	// Delete project users
	if err := tx.Where("project_id = ?", projectID).Delete(&models.UserProject{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update tasks to remove project association
	if err := tx.Model(&models.Task{}).Where("project_id = ?", projectID).Update("project_id", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update collaborative tasks to remove project association
	if err := tx.Model(&models.CollaborativeTask{}).Where("project_id = ?", projectID).Update("project_id", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the project
	if err := tx.Delete(&models.Project{}, projectID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetProjectTasks returns all tasks (regular and collaborative) for a project
func (s *ProjectService) GetProjectTasks(projectID uint) ([]models.Task, []models.CollaborativeTask, error) {
	var tasks []models.Task
	var collaborativeTasks []models.CollaborativeTask

	err := s.DB.Where("project_id = ?", projectID).Find(&tasks).Error
	if err != nil {
		return nil, nil, err
	}

	err = s.DB.Where("project_id = ?", projectID).Find(&collaborativeTasks).Error
	if err != nil {
		return nil, nil, err
	}

	return tasks, collaborativeTasks, nil
}

// GetProjectProgress returns the progress statistics for a project
func (s *ProjectService) GetProjectProgress(projectID uint) (map[string]interface{}, error) {
	tasks, collaborativeTasks, err := s.GetProjectTasks(projectID)
	if err != nil {
		return nil, err
	}

	var totalTasks, completedTasks, inProgressTasks, pendingTasks int

	// Count regular tasks
	for _, task := range tasks {
		totalTasks++
		switch task.Status {
		case models.TaskStatusCompleted:
			completedTasks++
		case models.TaskStatusInProgress:
			inProgressTasks++
		case models.TaskStatusPending:
			pendingTasks++
		}
	}

	// Count collaborative tasks
	for _, task := range collaborativeTasks {
		totalTasks++
		switch task.Status {
		case models.TaskStatusCompleted:
			completedTasks++
		case models.TaskStatusInProgress:
			inProgressTasks++
		case models.TaskStatusPending:
			pendingTasks++
		}
	}

	progress := 0.0
	if totalTasks > 0 {
		progress = float64(completedTasks) / float64(totalTasks) * 100
	}

	return map[string]interface{}{
		"total_tasks":       totalTasks,
		"completed_tasks":   completedTasks,
		"in_progress_tasks": inProgressTasks,
		"pending_tasks":     pendingTasks,
		"progress":          progress,
	}, nil
}
