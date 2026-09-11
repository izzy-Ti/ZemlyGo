package domain

import "time"

type Drivers struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserId     uint      `gorm:"column:user_id;not null;uniqueIndex" json:"user_id"`
	LicenseNo  string    `gorm:"not null" json:"license_no"`
	IsApproved bool      `gorm:"default:false" json:"is_approved"`
	IsOnline   bool      `gorm:"default:false" json:"is_online"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
