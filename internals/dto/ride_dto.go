package dto

import "time"

type EstimateFareRequest struct {
	PickupLat      float64 `json:"pickup_lat" binding:"required"`
	PickupLng      float64 `json:"pickup_lng" binding:"required"`
	DestinationLat float64 `json:"destination_lat" binding:"required"`
	DestinationLng float64 `json:"destination_lng" binding:"required"`
	RideType       string  `json:"ride_type"` // default standard
}

type FareOption struct {
	RideType   string  `json:"ride_type"`
	Fare       float64 `json:"fare"`
	DistanceKm float64 `json:"distance_km"`
	DurationMin int    `json:"estimated_duration_min"`
}

type EstimateFareResponse struct {
	DistanceKm  float64      `json:"distance_km"`
	Options     []FareOption `json:"options"`
}

type RequestRideRequest struct {
	PickupLat          float64 `json:"pickup_lat" binding:"required"`
	PickupLng          float64 `json:"pickup_lng" binding:"required"`
	PickupAddress      string  `json:"pickup_address"`
	DestinationLat     float64 `json:"destination_lat" binding:"required"`
	DestinationLng     float64 `json:"destination_lng" binding:"required"`
	DestinationAddress string  `json:"destination_address"`
	RideType           string  `json:"ride_type"` // standard, van, family
}

type CancelRideRequest struct {
	Reason string `json:"reason"`
}

type RideDTO struct {
	ID                 uint        `json:"id"`
	RiderID            uint        `json:"rider_id"`
	DriverID           *uint       `json:"driver_id,omitempty"`
	PickupLat          float64     `json:"pickup_lat"`
	PickupLng          float64     `json:"pickup_lng"`
	PickupAddress      string      `json:"pickup_address"`
	DestinationLat     float64     `json:"destination_lat"`
	DestinationLng     float64     `json:"destination_lng"`
	DestinationAddress string      `json:"destination_address"`
	RideType           string      `json:"ride_type"`
	Status             string      `json:"status"`
	Fare               float64     `json:"fare"`
	Rider              *UserDTO    `json:"rider,omitempty"`
	Driver             *DriverDTO  `json:"driver,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}
