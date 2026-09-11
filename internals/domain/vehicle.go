package domain

import "time"

type Vehicle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DriverID    uint      `gorm:"column:driver_id;not null;index" json:"driver_id"`
	PlateNumber string    `gorm:"column:plate_number;uniqueIndex;not null" json:"plate_number"`
	Model       string    `gorm:"not null" json:"model"`
	Color       string    `gorm:"not null" json:"color"`
	Status      bool      `gorm:"not null;default:true" json:"status"` // true = active vehicle for driver
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
