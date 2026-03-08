package domain

import "time"

type Rating struct {
	ID         uint `gorm:"primaryKey"`
	RideID     uint `gorm:"not null"`
	FromUserID uint `gorm:"not null"`
	ToUserID   uint `gorm:"not null"`
	Score      int  `gorm:"not null"`
	Comment    string
	CreatedAt  time.Time
}
