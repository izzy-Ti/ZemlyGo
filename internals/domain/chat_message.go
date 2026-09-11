package domain

import "time"

type ChatMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RideID     uint      `gorm:"column:ride_id;index;not null" json:"ride_id"`
	SenderID   uint      `gorm:"column:sender_id;not null" json:"sender_id"`
	SenderRole string    `gorm:"column:sender_role;not null" json:"sender_role"` // "rider" or "driver"
	Message    string    `gorm:"column:message;not null" json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
