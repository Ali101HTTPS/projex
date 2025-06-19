package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
	RoleManger   Role = "manger"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"not null"`
}

type Claims struct {
	UserID uint `json:"userId"`
	Role   Role `json:"role"`
	jwt.RegisteredClaims
}
