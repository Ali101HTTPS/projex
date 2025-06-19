package services

import (
	"errors"
	"project-x/models"
	"time"

	"gorm.io/gorm"
)

type CollaborativeTaskService struct {
	DB *gorm.DB
}

func NewCollaborativeTaskService(db *gorm.DB) *CollaborativeTaskService {
	return &CollaborativeTaskService{DB: db}
}

// CreateCollaborativeTask creates a new collaborative task with a lead user
func (s *CollaborativeTaskService) CreateCollaborativeTask(title, description string, leadUserID uint, projectID *uint, dueDate *time.Time, priority, complexity string) (*models.CollaborativeTask, error) {
	// Verify lead user exists
	var leadUser models.User
	if err := s.DB.First(&leadUser, leadUserID).Error; err != nil {
		return nil, errors.New("lead user not found")
	}

	// If projectID is provided, verify project exists and lead user is member
	if projectID != nil {
		var userProject models.UserProject
		if err := s.DB.Where("user_id = ? AND project_id = ?", leadUserID, *projectID).First(&userProject).Error; err != nil {
			return nil, errors.New("lead user is not a member of this project")
		}
	}

	// Validate priority
	validPriorities := []string{"high", "medium", "low"}
	validPriority := false
	for _, p := range validPriorities {
		if p == priority {
			validPriority = true
			break
		}
	}
	if !validPriority {
		return nil, errors.New("invalid priority")
	}

	// Validate complexity
	validComplexities := []string{"simple", "medium", "complex"}
	validComplexity := false
	for _, c := range validComplexities {
		if c == complexity {
			validComplexity = true
			break
		}
	}
	if !validComplexity {
		return nil, errors.New("invalid complexity")
	}

	task := &models.CollaborativeTask{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		LeadUserID:  leadUserID,
		ProjectID:   projectID,
		AssignedAt:  time.Now(),
		DueDate:     dueDate,
		Priority:    priority,
		Complexity:  complexity,
		Progress:    0,
	}

	if err := s.DB.Create(task).Error; err != nil {
		return nil, err
	}

	// Automatically add lead user as participant with lead role
	participant := &models.CollaborativeTaskParticipant{
		CollaborativeTaskID: task.ID,
		UserID:              leadUserID,
		Role:                "lead",
		Status:              "active",
		AssignedAt:          time.Now(),
	}

	if err := s.DB.Create(participant).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// AddParticipant adds a user to a collaborative task
func (s *CollaborativeTaskService) AddParticipant(taskID, userID uint, role, contribution string) error {
	// Verify task exists
	var task models.CollaborativeTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return errors.New("collaborative task not found")
	}

	// Verify user exists
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if user is already a participant
	var existingParticipant models.CollaborativeTaskParticipant
	if err := s.DB.Where("collaborative_task_id = ? AND user_id = ?", taskID, userID).First(&existingParticipant).Error; err == nil {
		return errors.New("user is already a participant in this task")
	}

	// Validate role
	validRoles := []string{"lead", "contributor", "reviewer", "observer"}
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

	participant := &models.CollaborativeTaskParticipant{
		CollaborativeTaskID: taskID,
		UserID:              userID,
		Role:                role,
		Status:              "active",
		AssignedAt:          time.Now(),
		Contribution:        contribution,
	}

	return s.DB.Create(participant).Error
}

// RemoveParticipant removes a user from a collaborative task
func (s *CollaborativeTaskService) RemoveParticipant(taskID, userID uint) error {
	// Check if user is the lead (cannot remove lead)
	var task models.CollaborativeTask
	if err := s.DB.First(&task, taskID).Error; err != nil {
		return errors.New("collaborative task not found")
	}

	if task.LeadUserID == userID {
		return errors.New("cannot remove the lead user from the task")
	}

	result := s.DB.Where("collaborative_task_id = ? AND user_id = ?", taskID, userID).Delete(&models.CollaborativeTaskParticipant{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user is not a participant in this task")
	}

	return nil
}

// UpdateTaskProgress updates the progress of a collaborative task
func (s *CollaborativeTaskService) UpdateTaskProgress(taskID uint, progress int) error {
	if progress < 0 || progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}

	// Update progress
	if err := s.DB.Model(&models.CollaborativeTask{}).Where("id = ?", taskID).Update("progress", progress).Error; err != nil {
		return err
	}

	// If progress is 100%, mark task as completed
	if progress == 100 {
		if err := s.DB.Model(&models.CollaborativeTask{}).Where("id = ?", taskID).Update("status", models.TaskStatusCompleted).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetCollaborativeTaskWithDetails returns a collaborative task with all related data
func (s *CollaborativeTaskService) GetCollaborativeTaskWithDetails(taskID uint) (*models.CollaborativeTask, error) {
	var task models.CollaborativeTask
	err := s.DB.Preload("LeadUser").
		Preload("Project").
		Preload("Participants.User").
		First(&task, taskID).Error
	return &task, err
}

// GetUserCollaborativeTasks returns all collaborative tasks where user is a participant
func (s *CollaborativeTaskService) GetUserCollaborativeTasks(userID uint) ([]models.CollaborativeTask, error) {
	var tasks []models.CollaborativeTask
	err := s.DB.Joins("JOIN collaborative_task_participants ON collaborative_tasks.id = collaborative_task_participants.collaborative_task_id").
		Where("collaborative_task_participants.user_id = ?", userID).
		Preload("LeadUser").
		Preload("Project").
		Preload("Participants.User").
		Order("collaborative_tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// GetCollaborativeTaskStatistics returns statistics for a collaborative task
func (s *CollaborativeTaskService) GetCollaborativeTaskStatistics(taskID uint) (map[string]interface{}, error) {
	var task models.CollaborativeTask
	if err := s.DB.Preload("Participants").First(&task, taskID).Error; err != nil {
		return nil, errors.New("collaborative task not found")
	}

	// Count participants by role
	var leadCount, contributorCount, reviewerCount, observerCount int
	for _, participant := range task.Participants {
		switch participant.Role {
		case "lead":
			leadCount++
		case "contributor":
			contributorCount++
		case "reviewer":
			reviewerCount++
		case "observer":
			observerCount++
		}
	}

	stats := map[string]interface{}{
		"task_id":            task.ID,
		"title":              task.Title,
		"status":             task.Status,
		"progress":           task.Progress,
		"priority":           task.Priority,
		"complexity":         task.Complexity,
		"total_participants": len(task.Participants),
		"participants_by_role": map[string]int{
			"lead":        leadCount,
			"contributor": contributorCount,
			"reviewer":    reviewerCount,
			"observer":    observerCount,
		},
		"assigned_at": task.AssignedAt,
		"due_date":    task.DueDate,
	}

	return stats, nil
}
