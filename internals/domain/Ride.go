package domain

import "time"

type Ride struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	RiderID            uint       `gorm:"column:rider_id;not null;index" json:"rider_id"`
	DriverID           *uint      `gorm:"column:driver_id;index" json:"driver_id,omitempty"`
	PickupLat          float64    `gorm:"column:pickup_lat;not null" json:"pickup_lat"`
	PickupLng          float64    `gorm:"column:pickup_lng;not null" json:"pickup_lng"`
	PickupAddress      string     `gorm:"column:pickup_address" json:"pickup_address"`
	DestinationLat     float64    `gorm:"column:destination_lat;not null" json:"destination_lat"`
	DestinationLng     float64    `gorm:"column:destination_lng;not null" json:"destination_lng"`
	DestinationAddress string     `gorm:"column:destination_address" json:"destination_address"`
	RideType           string     `gorm:"column:ride_type;not null;default:standard" json:"ride_type"`
	Status             string     `gorm:"column:status;not null;default:requested;index" json:"status"`
	Fare               float64    `gorm:"column:fare;default:0" json:"fare"`
	ScheduledAt        *time.Time `gorm:"column:scheduled_at" json:"scheduled_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
