package models

import (
	"time"
)

type Organization struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"not null" json:"name"`
    Description string    `json:"description"`
    Address     string    `json:"address"`
    City        string    `json:"city"`
    State       string    `json:"state"`
    Country     string    `json:"country"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
    DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}
