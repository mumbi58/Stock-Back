package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"not null" json:"password" validate:"required,min=6"`
	FirstName string         `gorm:"not null" json:"first_name" validate:"required,min=2,max=50"`
	LastName  string         `gorm:"not null" json:"last_name" validate:"required,min=2,max=50"`
	RoleID    uint           `gorm:"not null" json:"role_id" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}
