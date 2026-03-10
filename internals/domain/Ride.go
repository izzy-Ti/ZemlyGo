package domain

import "time"

type Ride struct {
	ID                 uint `gorm:"primaryKey"`
	RiderID            uint `gorm:"not null"`
	DriverID           *uint
	PickupLat          float64 `gorm:"not null"`
	PickupLng          float64 `gorm:"not null"`
	PickupAddress      string
	DestinationLat     float64 `gorm:"not null"`
	DestinationLng     float64 `gorm:"not null"`
	DestinationAddress string
	RideType           string  `gorm:"not null;default:standard"`
	Status             string  `gorm:"not null;default:requested"`
	Fare               float64 `gorm:"default:0"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
