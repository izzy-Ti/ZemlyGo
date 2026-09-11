package dto

import "time"

type RegisterDriverRequest struct {
	LicenseNo string `json:"license_no" binding:"required"`
}

type DriverDTO struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	LicenseNo  string     `json:"license_no"`
	IsApproved bool       `json:"is_approved"`
	IsOnline   bool       `json:"is_online"`
	User       *UserDTO   `json:"user,omitempty"`
	Vehicle    *VehicleDTO `json:"active_vehicle,omitempty"`
	Location   *LocationDTO `json:"location,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type UpdateLocationRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}

type LocationDTO struct {
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NearbyDriverDTO struct {
	DriverID   uint        `json:"driver_id"`
	DistanceKm float64     `json:"distance_km"`
	Lat        float64     `json:"lat"`
	Lng        float64     `json:"lng"`
	Vehicle    *VehicleDTO `json:"vehicle,omitempty"`
}
