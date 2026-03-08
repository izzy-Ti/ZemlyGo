package domain

import "time"

type Vehcile struct {
	ID          uint   `gorm:"primaryKey"`
	DriverID    uint   `gorm:"not null"`
	PlateNumber string `gorm:"uniqueIndex;not null"`
	Model       string `gorm:"not null"`
	Color       string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
