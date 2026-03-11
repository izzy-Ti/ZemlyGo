package domain

import "time"

type DriverLocation struct {
	ID        uint    `gorm:"primaryKey"`
	DriverID  uint    `gorm:"index;not null"`
	Lat       float64 `gorm:"not null"`
	Lng       float64 `gorm:"not null"`
	UpdatedAt time.Time
}
