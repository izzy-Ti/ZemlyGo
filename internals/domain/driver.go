package domain

import "time"

type Drivers struct {
	ID         uint   `gorm:"primaryKey"`
	UserId     uint   `gorm:"not null"`
	LicenseNo  string `gorm:"not null"`
	IsApproved bool   `gorm:"default:false"`
	IsOnline   bool   `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
