package services

import (
	"errors"
	"project-x/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(username, password, role, department string) (*models.User, error) {
	// Validate role
	if !s.isValidRole(role) {
		return nil, errors.New("invalid role")
	}

	// Check if username already exists
	var existingUser models.User
	if err := s.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:   username,
		Password:   string(hashedPassword),
		Role:       models.Role(role),
		Department: department,
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername returns a user by username
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers returns all users with optional filtering
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := s.DB.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByRole returns users filtered by role
func (s *UserService) GetUsersByRole(role string) ([]models.User, error) {
	if !s.isValidRole(role) {
		return nil, errors.New("invalid role")
	}

	var users []models.User
	if err := s.DB.Where("role = ?", role).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByDepartment returns users filtered by department
func (s *UserService) GetUsersByDepartment(department string) ([]models.User, error) {
	var users []models.User
	if err := s.DB.Where("department = ?", department).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(userID uint, newRole string) (*models.User, error) {
	if !s.isValidRole(newRole) {
		return nil, errors.New("invalid role")
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = models.Role(newRole)
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserDepartment updates a user's department
func (s *UserService) UpdateUserDepartment(userID uint, newDepartment string) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Department = newDepartment
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserPassword updates a user's password
func (s *UserService) UpdateUserPassword(userID uint, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error
}

// DeleteUser deletes a user and all related data
func (s *UserService) DeleteUser(userID uint) error {
	// Start a transaction
	tx := s.DB.Begin()

	// Check if user exists
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return errors.New("user not found")
	}

	// Delete user project memberships
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserProject{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user's tasks (they will be deleted due to CASCADE constraint)
	if err := tx.Where("user_id = ?", userID).Delete(&models.Task{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user's collaborative tasks (they will be deleted due to CASCADE constraint)
	if err := tx.Where("user_id = ?", userID).Delete(&models.CollaborativeTask{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserStats returns statistics about a user
func (s *UserService) GetUserStats(userID uint) (map[string]interface{}, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Count user's tasks
	var taskCount int64
	s.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&taskCount)

	// Count user's collaborative tasks
	var collaborativeTaskCount int64
	s.DB.Model(&models.CollaborativeTask{}).Where("user_id = ?", userID).Count(&collaborativeTaskCount)

	// Count user's projects
	var projectCount int64
	s.DB.Model(&models.UserProject{}).Where("user_id = ?", userID).Count(&projectCount)

	// Count completed tasks
	var completedTaskCount int64
	s.DB.Model(&models.Task{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusCompleted).Count(&completedTaskCount)

	// Count completed collaborative tasks
	var completedCollaborativeTaskCount int64
	s.DB.Model(&models.CollaborativeTask{}).Where("user_id = ? AND status = ?", userID, models.TaskStatusCompleted).Count(&completedCollaborativeTaskCount)

	return map[string]interface{}{
		"user_id":                       userID,
		"username":                      user.Username,
		"role":                          user.Role,
		"department":                    user.Department,
		"total_tasks":                   taskCount,
		"total_collaborative_tasks":     collaborativeTaskCount,
		"total_projects":                projectCount,
		"completed_tasks":               completedTaskCount,
		"completed_collaborative_tasks": completedCollaborativeTaskCount,
		"completion_rate":               s.calculateCompletionRate(taskCount+collaborativeTaskCount, completedTaskCount+completedCollaborativeTaskCount),
		"created_at":                    user.CreatedAt,
		"last_active":                   user.UpdatedAt,
	}, nil
}

// isValidRole checks if a role is valid
func (s *UserService) isValidRole(role string) bool {
	validRoles := []string{"admin", "head", "manager", "employee"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// calculateCompletionRate calculates the completion rate percentage
func (s *UserService) calculateCompletionRate(total, completed int64) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(completed) / float64(total) * 100
}
