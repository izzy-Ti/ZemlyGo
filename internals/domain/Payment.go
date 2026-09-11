package domain

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RideID    uint      `gorm:"column:ride_id;not null;uniqueIndex" json:"ride_id"`
	Amount    float64   `gorm:"column:amount;not null" json:"amount"`
	Tip       float64   `gorm:"column:tip;default:0" json:"tip"`
	Method    string    `gorm:"column:method;not null" json:"method"` // card, cash, wallet
	Status    string    `gorm:"column:status;not null;default:pending" json:"status"` // pending, completed, failed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
