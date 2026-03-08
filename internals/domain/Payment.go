package domain

import "time"

type Payment struct {
	ID        uint    `gorm:"primaryKey"`
	RideID    uint    `gorm:"not null"`
	Amount    float64 `gorm:"not null"`
	Method    string  `gorm:"not null"`
	Status    string  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
