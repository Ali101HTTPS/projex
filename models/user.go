package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
	RoleManager  Role = "manager"
	RoleHead     Role = "head"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusPaused    ProjectStatus = "paused"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

type User struct {
	gorm.Model
	Username   string `gorm:"unique;not null;index"`
	Password   string `gorm:"not null"`
	Role       Role   `gorm:"not null;index"`
	Department string `gorm:"not null;index"`

	// Relationships
	Tasks              []Task              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CollaborativeTasks []CollaborativeTask `gorm:"foreignKey:LeadUserID;constraint:OnDelete:CASCADE"`   // Tasks where user is the lead
	Projects           []Project           `gorm:"many2many:user_projects;constraint:OnDelete:CASCADE"` // Many-to-many relationship
	CreatedProjects    []Project           `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL"`   // Projects created by this user

	// Collaborative task participation
	CollaborativeTaskParticipations []CollaborativeTaskParticipant `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Task struct {
	gorm.Model
	Title       string     `gorm:"not null;index"`
	Description string     `gorm:"not null"`
	Status      TaskStatus `gorm:"not null;default:'pending';index"`
	UserID      uint       `gorm:"not null;index"`
	ProjectID   *uint      `gorm:"index"` // Optional: task can belong to a project
	AssignedAt  time.Time  `gorm:"not null;index"`
	DueDate     *time.Time `gorm:"index"` // Optional due date

	// Relationships
	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Project *Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
}

type CollaborativeTask struct {
	gorm.Model
	Title       string     `gorm:"not null;index"`
	Description string     `gorm:"not null"`
	Status      TaskStatus `gorm:"not null;default:'pending';index"`
	LeadUserID  uint       `gorm:"not null;index"` // Main person responsible
	ProjectID   *uint      `gorm:"index"`          // Optional: collaborative task can belong to a project
	AssignedAt  time.Time  `gorm:"not null;index"`
	DueDate     *time.Time `gorm:"index"`                  // Optional due date
	Priority    string     `gorm:"default:'medium';index"` // high, medium, low
	Progress    int        `gorm:"default:0;index"`        // 0-100 percentage
	Complexity  string     `gorm:"default:'medium';index"` // simple, medium, complex

	// Relationships
	LeadUser     User                           `gorm:"foreignKey:LeadUserID;constraint:OnDelete:SET NULL"`
	Project      *Project                       `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Participants []CollaborativeTaskParticipant `gorm:"foreignKey:CollaborativeTaskID;constraint:OnDelete:CASCADE"`
}

// CollaborativeTaskParticipant represents team members working on a collaborative task
type CollaborativeTaskParticipant struct {
	gorm.Model
	CollaborativeTaskID uint       `gorm:"not null;index"`
	UserID              uint       `gorm:"not null;index"`
	Role                string     `gorm:"not null;default:'contributor';index"` // lead, contributor, reviewer, observer
	Status              string     `gorm:"not null;default:'active';index"`      // active, inactive, completed
	AssignedAt          time.Time  `gorm:"not null;index"`
	CompletedAt         *time.Time `gorm:"index"`
	Contribution        string     `gorm:"index"` // Description of their contribution

	// Relationships
	CollaborativeTask CollaborativeTask `gorm:"foreignKey:CollaborativeTaskID;constraint:OnDelete:CASCADE"`
	User              User              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Project struct {
	gorm.Model
	Title       string        `gorm:"not null;index"`
	Description string        `gorm:"not null"`
	Status      ProjectStatus `gorm:"not null;default:'active';index"`
	CreatedBy   uint          `gorm:"not null;index"` // Who created the project
	StartDate   time.Time     `gorm:"not null;index"`
	EndDate     *time.Time    `gorm:"index"` // Optional end date

	// Relationships
	Creator            User                `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL"`
	Users              []User              `gorm:"many2many:user_projects;constraint:OnDelete:CASCADE"` // Many-to-many relationship
	Tasks              []Task              `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	CollaborativeTasks []CollaborativeTask `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

// UserProject represents the many-to-many relationship between users and projects
type UserProject struct {
	UserID    uint      `gorm:"primaryKey;index"`
	ProjectID uint      `gorm:"primaryKey;index"`
	JoinedAt  time.Time `gorm:"not null;index"`
	Role      string    `gorm:"not null;default:'member';index"` // member, lead, contributor, etc.

	// Relationships
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

type Claims struct {
	UserID uint `json:"userId"`
	Role   Role `json:"role"`
	jwt.RegisteredClaims
}
