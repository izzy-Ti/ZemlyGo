package domain

import "time"

type EmergencyAlert struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RideID    uint      `gorm:"column:ride_id;index;not null" json:"ride_id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	Lat       float64   `gorm:"column:lat;not null" json:"lat"`
	Lng       float64   `gorm:"column:lng;not null" json:"lng"`
	Reason    string    `gorm:"column:reason" json:"reason"`
	Status    string    `gorm:"column:status;default:open" json:"status"` // "open", "dispatched", "resolved"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
