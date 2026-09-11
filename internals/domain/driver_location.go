package domain

import "time"

type DriverLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DriverID  uint      `gorm:"column:driver_id;uniqueIndex;not null" json:"driver_id"`
	Lat       float64   `gorm:"column:lat;not null" json:"lat"`
	Lng       float64   `gorm:"column:lng;not null" json:"lng"`
	UpdatedAt time.Time `json:"updated_at"`
}
