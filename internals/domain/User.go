package domain

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	SupabaseUserID string `gorm:"uniqueIndex;not null"`
	Name           string `gorm:"not null"`
	Email          string `gorm:"uniqueIndex;not null"`
	Phone          string `gorm:"not null"`
	Role           string `gorm:"not null;default:rider"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
