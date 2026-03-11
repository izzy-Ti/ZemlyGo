package domain

import "time"

type Vehicle struct {
	ID          uint   `gorm:"primaryKey"`
	DriverID    uint   `gorm:"not null"`
	PlateNumber string `gorm:"uniqueIndex;not null"`
	Model       string `gorm:"not null"`
	Color       string `gorm:"not null"`
	Status      bool   `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
